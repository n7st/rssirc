// Package rssfeed runs in parallel to rssbot. It polls an RSS feed URL and
// sends formatted output to the IRC bot.
//
// This file handles polling of an individual RSS feed.
package rssfeed

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"

	"github.com/n7st/rssirc/internal/app/rssbot"
	"github.com/n7st/rssirc/internal/pkg/util"
)

// Poller contains details for polling RSS feeds.
type Poller struct {
	Bot    *rssbot.Bot
	Cache  *Cache
	Config *util.RSSConfig
	Logger *logrus.Logger
	Parser *gofeed.Parser

	// If true, the poller is on its first run and shouldn't announce the last
	// items from the feed.
	FirstRun bool
}

// NewPoller sets up a new web poller for scraping the given RSS feed.
func NewPoller(config *util.RSSConfig, bot *rssbot.Bot, logger *logrus.Logger) *Poller {
	return &Poller{
		Bot:      bot,
		Cache:    NewCache(config.MaxHistory),
		Config:   config,
		FirstRun: true,
		Logger:   logger,
		Parser:   gofeed.NewParser(),
	}
}

// Poll polls the feed for new items. The user-provided delay is converted into
// minutes, and the poller waits that long between tries.
func (p *Poller) Poll() {
	for {
		// Allow time for connection to IRC to be made
		if p.Bot.Connection.Connected() {
			// Allow time for the bot to join channels
			time.Sleep(1 * time.Second)
			break
		} else {
			time.Sleep(10 * time.Second)
		}
	}

	go func() {
		for {
			p.Logger.WithFields(logrus.Fields{
				"url": p.Config.FeedURL,
			}).Debug("Polling")

			feed, err := p.Parser.ParseURL(p.Config.FeedURL)

			if err == nil {
				p.announce(feed)
			} else {
				p.Logger.WithFields(logrus.Fields{
					"error": err.Error(),
				}).Warn("An error occurred polling the RSS feed")
			}

			time.Sleep(p.Config.PollDelayMinutes)
		}
	}()
}

// announce puts feed notifications in IRC channels.
func (p *Poller) announce(feed *gofeed.Feed) {
	for _, item := range feed.Items[0:p.Config.MaxHistory] {
		if !p.Cache.Exists(item.Title) {
			p.Cache.Save(item)

			if !p.FirstRun {
				message := fmt.Sprintf("%s %s", item.Title, item.Link)

				if len(p.Config.Channels) > 0 {
					p.Bot.MessageChannels(p.Config.Channels, message)
				} else {
					p.Logger.WithFields(logrus.Fields{
						"message": message,
						"url":     p.Config.FeedURL,
					}).Warn("No channels provided for message")
				}
			}
		}
	}

	if p.FirstRun {
		// The first run should just populate the cache. Future runs
		// will announce new feed items.
		p.Logger.Debug("Cache populated")
		p.FirstRun = false
	}
}
