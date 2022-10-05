package main

import (
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/exceptioon/tiktok-fav-publisher/config"
	"github.com/exceptioon/tiktok-fav-publisher/internal"
	"github.com/exceptioon/tiktok-fav-publisher/internal/store/redis"
	"github.com/exceptioon/tiktok-fav-publisher/internal/store/set"
	"github.com/exceptioon/tiktok-fav-publisher/internal/tiktok"
	"github.com/exceptioon/tiktok-fav-publisher/internal/worker"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
)

func main() {
	var (
		conf     config.Config
		bot      worker.Bot
		cache    internal.Cache
		wg       sync.WaitGroup
		tt       *tiktok.ServiceApi
		logger   *zap.Logger
		tick     = time.NewTicker(time.Minute)
		quitChan chan struct{}
		err      error
	)

	err = env.Parse(&conf)
	panicErr(err)

	bot.Bot, err = telebot.NewBot(telebot.Settings{
		Token:  conf.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: time.Second * 10},
	})
	bot.ChatID = conf.ChannelID

	switch conf.CacheType {
	case config.CacheTypeRedis:
		cache, err = redis.NewRedisCache(conf.DBAddr)
		panicErr(err)
	case config.CacheTypeSet:
		cache, err = set.NewSet()
		panicErr(err)
	default:
		panicErr(errors.New("not realized yet"))
	}

	logger, err = zap.NewProduction()
	panicErr(err)

	tt = tiktok.NewServiceApi(conf.TikTokUsername, 10, logger)

	w := worker.Worker{
		TikTok:   tt,
		TG:       bot,
		Cache:    cache,
		WG:       &wg,
		Tick:     tick,
		Log:      logger,
		QuitChan: quitChan,
	}

	wg.Add(1)
	go func() {
		w.Start()
	}()

	logger.Info("worker started")

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	sig := <-shutdownSignal
	logger.Info("shutting down...", zap.String("signal", sig.String()))
	w.Stop()
	wg.Wait()

	close(shutdownSignal)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
