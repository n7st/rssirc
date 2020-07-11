// Package rssbot is a simple IRC bot which connects to a single IRC network and
// receives messages from the rssfeed package.
//
// This file sets up an IRC bot instance.
package rssbot

import (
	"crypto/tls"
	"time"

	"github.com/sirupsen/logrus"
	irc "github.com/thoj/go-ircevent"

	"gitlab.com/n7st/rssirc/internal/pkg/util"
)

// Bot contains the IRC bot.
type Bot struct {
	Connection *irc.Connection
	Config     *util.Config
	Logger     *logrus.Logger
}

// Init sets up an IRC bot connection to the network.
func Init(config *util.Config, logger *logrus.Logger) *Bot {
	connection := irc.IRC(config.IRC.Nickname, config.IRC.Ident)

	connection.Debug = config.IRC.Debug
	connection.VerboseCallbackHandler = config.IRC.Verbose
	connection.RealName = config.IRC.RealName

	if config.IRC.UseTLS {
		connection.UseTLS = true
		connection.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	err := connection.Connect(config.IRC.Hostname)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Fatal error connecting to IRC")
	}

	bot := &Bot{Connection: connection, Config: config, Logger: logger}

	for name, fn := range events(bot) {
		bot.Connection.AddCallback(name, fn)
	}

	bot.healthCheck()

	return bot
}

// joinChannels joins either a provided list of channels, or the channels set in
// the bot's configuration.
func (b *Bot) joinChannels(params ...string) {
	var channels []string

	if len(params) > 0 {
		channels = params
	} else {
		channels = b.Config.IRC.Channels
	}

	for _, channel := range channels {
		b.Connection.Join(channel)
	}
}

// MessageChannels sends one message to many channels.
func (b *Bot) MessageChannels(channels []string, message string) {
	for _, channel := range channels {
		b.Connection.Privmsg(channel, message)

		time.Sleep(1 * time.Second) // Antispam
	}
}

// healthCheck checks for a connection to the IRC network and reconnects as
// required.
func (b *Bot) healthCheck() {
	retries := 0

	go func() {
		for {
			select {
			case <-b.Connection.Error:
				b.Logger.Warn("Healthcheck failed")

				if retries > b.Config.IRC.MaxReconnect {
					b.Logger.WithFields(logrus.Fields{
						"retries":     retries,
						"max_retries": b.Config.IRC.MaxReconnect,
					}).Fatal("Maximum reconnection attempts exceeded")
				}

				err := b.Connection.Reconnect()

				if err != nil {
					b.Logger.WithFields(logrus.Fields{
						"error": err.Error(),
					}).Warn("Health check error")
				} else {
					retries = 0
				}
			default:
				if b.Config.IRC.Verbose {
					b.Logger.Debug("Health check successful")
				}
			}

			time.Sleep(b.Config.IRC.ReconnectDelayMinutes * time.Minute)
		}
	}()
}
