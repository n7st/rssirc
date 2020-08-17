// Package main contains the application's entrypoint.
package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/n7st/rssirc/internal/app/rssbot"
	"github.com/n7st/rssirc/internal/app/rssfeed"
	"github.com/n7st/rssirc/internal/pkg/util"
)

// main sets up an IRC bot and many RSS feed pollers.
func main() {
	var config *util.Config

	if len(os.Args) > 1 {
		config = util.NewConfig(os.Args[1])
	} else {
		config = util.NewConfig()
	}

	logger := logrus.New()

	logger.SetLevel(config.LogLevel)
	logger.WithFields(logrus.Fields{
		"level": config.LogLevel.String(),
	}).Info("Set log level")

	bot := rssbot.Init(config, logger)

	for _, rss := range config.RSS {
		poller := rssfeed.NewPoller(rss, bot, logger)
		poller.Poll()
	}

	bot.Connection.Loop()
}
