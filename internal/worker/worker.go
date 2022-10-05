package worker

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/exceptioon/tiktok-fav-publisher/internal"
	"github.com/exceptioon/tiktok-fav-publisher/internal/tiktok"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	Bot    *telebot.Bot
	ChatID int64
}

type Worker struct {
	TikTok *tiktok.ServiceApi
	TG     Bot
	Cache  internal.Database
	WG     *sync.WaitGroup
	Log    *zap.Logger
	Tick   *time.Ticker

	QuitChan chan struct{}
}

func (w *Worker) Start() {
	defer w.WG.Done()

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
				time.Sleep(time.Second * 10)
			}

		case <-w.QuitChan:
			w.Tick.Stop()
			break LOOP
		}
	}
}

func (w *Worker) Stop() {
	w.QuitChan <- struct{}{}
}
