package worker

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/exceptioon/tiktok-fav-publisher/internal"
	"github.com/exceptioon/tiktok-fav-publisher/internal/tiktok"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

var tiktokVideo = regexp.MustCompile(`(?m)(vm|www)\.tiktok\.com/.+`)

type Bot struct {
	Bot    *telebot.Bot
	ChatID int64
}

type Worker struct {
	TikTok *tiktok.ServiceApi
	TG     Bot
	Cache  internal.Cache
	WG     *sync.WaitGroup
	Log    *zap.Logger
	Tick   *time.Ticker

	QuitChan chan struct{}
}

func (w *Worker) Start() {
	defer w.WG.Done()
	go w.HandleTelegramQueries()

LOOP:
	for {
		select {
		case <-w.Tick.C:
			videos, err := w.TikTok.GetLikedVideos()
			if err != nil {
				w.Log.Error("got error", zap.Error(err))
				continue
			}
			for _, video := range videos {
				if w.Cache.IsExist(video.ID) {
					continue
				}
				time.Sleep(time.Second * 5)
				err = w.TikTok.SetVideoMetadata(&video)
				w.Log.Info("processing video", zap.String("ID", video.ID), zap.String("url", video.ShareableLink),
					zap.String("download", video.DownloadLink))

				if strings.HasSuffix(video.DownloadLink, ".mp3") {
					err = w.Cache.Add(video.ID)
					if err != nil {
						w.Log.Error("Add to cache", zap.Error(err))
					}
					continue
				}

				menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
				menu.Inline(
					menu.Row(menu.URL("Original", video.ShareableLink)),
				)
				_, err = w.TG.Bot.Send(telebot.ChatID(w.TG.ChatID),
					&telebot.Video{
						File:    telebot.File{FileURL: video.DownloadLink},
						Caption: fmt.Sprintf("@%s: %s", video.AuthorUsername, video.Title),
					}, menu)

				if err != nil {
					w.Log.Error("Send video", zap.Error(err), zap.String("download url", video.DownloadLink))
					continue
				}

				err = w.Cache.Add(video.ID)
				if err != nil {
					w.Log.Error("Add to cache", zap.Error(err))
					continue
				}
				w.Log.Info("sent video", zap.String("id", video.ID))
			}

		case <-w.QuitChan:
			w.Tick.Stop()
			w.TG.Bot.Stop()
			break LOOP
		}
	}
}

func (w *Worker) Stop() {
	w.QuitChan <- struct{}{}
}

func (w *Worker) HandleTelegramQueries() {
	w.TG.Bot.Handle(telebot.OnQuery, func(c telebot.Context) error {
		var (
			video tiktok.Video
			query = c.Query()
		)
		video.ID = query.Text
		if !tiktokVideo.Match([]byte(video.ID)) {
			return nil
		}
		err := w.TikTok.SetVideoMetadata(&video)
		if err != nil {
			w.Log.Error("HandleTelegramQueries SetVideoMetadata problem", zap.Error(err))
			return err
		}
		w.Log.Info("HandleTelegramQueries processing", zap.String("req", video.ID), zap.String("User", query.Sender.Username),
			zap.String("Name", query.Sender.FirstName+" "+query.Sender.LastName))

		results := telebot.Results{
			&telebot.VideoResult{
				ThumbURL: video.Cover,
				URL:      video.DownloadLink,
				Caption:  fmt.Sprintf("@%s: %s", video.AuthorUsername, video.Title),
				MIME:     "video/mp4",
				Title:    fmt.Sprintf("@%s: %s", video.AuthorUsername, video.Title),
			},
		}
		results[0].SetResultID("0")

		return c.Answer(&telebot.QueryResponse{
			Results:   results,
			CacheTime: 60,
		})
	})

	w.TG.Bot.Start()
}
