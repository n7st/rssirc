// Package rssbot is a simple IRC bot which connects to a single IRC network and
// receives messages from the rssfeed package.
//
// This file contains callback handlers for IRC numeric events.
package rssbot

import irc "github.com/thoj/go-ircevent"

// events returns all the bot's available IRC events.
func events(b *Bot) map[string]func(e *irc.Event) {
	return map[string]func(e *irc.Event){
		"001": b.callback001,
		"900": b.callback900,
	}
}

// callback001 runs when the bot connects to the network.
func (b *Bot) callback001(e *irc.Event) {
	if b.Config.IRC.Modes != "" {
		b.Connection.Mode(b.Connection.GetNick(), b.Config.IRC.Modes)
	}

	if b.Config.IRC.NickservPassword != "" {
		b.Connection.Privmsgf("nickserv", "identify %s", b.Config.IRC.NickservPassword)
	} else {
		b.joinChannels()
	}
}

// callback900 runs when the bot receives confirmation of nickserv login.
func (b *Bot) callback900(e *irc.Event) {
	b.joinChannels()
}
