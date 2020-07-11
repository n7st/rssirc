// Package util contains common functionality for "utilities" required by the
// bot and feed poller.
//
// This file is for processing the application's configuration from a YAML file.
package util

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/shibukawa/configdir"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	configFilename  = "config.yaml" // These default values are used to create a
	vendorName      = "netsplit"    // filename for loading a file from the
	applicationName = "rssirc"      // platform's standard config location.

	defaultCacheLen = 3
	defaultPort     = 6667
	defaultNickname = "rssirc"
	defaultLogLevel = "info"
)

// ircConfig contains config items specific to the IRC bot itself.
type ircConfig struct {
	Channels         []string `yaml:"channels"`
	Debug            bool     `yaml:"debug"`
	Ident            string   `yaml:"ident"`
	MaxReconnect     int      `yaml:"max_reconnect"`
	ReconnectDelay   int      `yaml:"reconnect_delay"`
	Modes            string   `yaml:"modes"`
	Nickname         string   `yaml:"nickname"`
	NickservPassword string   `yaml:"nickserv_password"`
	Port             int      `yaml:"port"`
	RealName         string   `yaml:"real_name"`
	Server           string   `yaml:"server"`
	ServerPassword   string   `yaml:"server_password"`
	UseTLS           bool     `yaml:"use_tls"`
	Verbose          bool     `yaml:"verbose"`

	Hostname              string
	ReconnectDelayMinutes time.Duration
}

// RSSConfig contains config items specific to the RSS feed poller.
type RSSConfig struct {
	Channels   []string `yaml:"channels"`
	FeedURL    string   `yaml:"feed_url"`
	PollDelay  int      `yaml:"poll_delay"`
	MaxHistory int      `yaml:"max_history"`

	PollDelayMinutes time.Duration
}

// Config contains the entire application's configuration.
type Config struct {
	IRC              *ircConfig   `yaml:"irc"`
	UnparsedLogLevel string       `yaml:"log_level"`
	RSS              []*RSSConfig `yaml:"rss"`

	LogLevel logrus.Level
}

// NewConfig sets up the application's configuration.
func NewConfig(params ...string) *Config {
	config := &Config{}

	data, err := loadConfigData(params)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		panic(err)
	}

	config.applyDefaults()

	return config
}

// loadConfigData retrieves bytes from the config file. An optional filename can
// be provided to load configuration from a specific file rather than the
// default set for configdir.
func loadConfigData(params []string) ([]byte, error) {
	var (
		err      error
		data     []byte
		filename string
	)

	if len(params) > 0 {
		filename = params[0]

		data, err = ioutil.ReadFile(filename)
	} else {
		configDirs := configdir.New(vendorName, applicationName)
		folder := configDirs.QueryFolderContainsFile(configFilename)

		if folder != nil {
			data, err = folder.ReadFile(configFilename)
		}
	}

	if err != nil {
		panic(err)
	}

	return data, err
}

// applyDefaults sets default configuration values for items which are missing.
func (c *Config) applyDefaults() {
	if c.IRC.Port == 0 {
		c.IRC.Port = defaultPort
	}

	if c.IRC.Nickname == "" {
		c.IRC.Nickname = defaultNickname
	}

	if c.IRC.Ident == "" {
		c.IRC.Ident = c.IRC.Nickname
	}

	if c.IRC.RealName == "" {
		c.IRC.RealName = c.IRC.Nickname
	}

	if c.IRC.ReconnectDelay == 0 {
		c.IRC.ReconnectDelay = 10
	}

	for _, poller := range c.RSS {
		if poller.MaxHistory == 0 {
			poller.MaxHistory = defaultCacheLen
		}

		if poller.PollDelay < 1 {
			panic("the minimum poll delay is 1 minute")
		}

		poller.PollDelayMinutes = time.Minute * intToDuration(poller.PollDelay)
	}

	c.IRC.Hostname = fmt.Sprintf("%s:%d", c.IRC.Server, c.IRC.Port)
	c.IRC.ReconnectDelayMinutes = intToDuration(c.IRC.ReconnectDelay)

	c.setLogLevel()
}

// intToDuration converts an integer into a time.Duration for use with
// time.Sleep.
func intToDuration(input int) time.Duration {
	return time.Duration(input)
}

// setLogLevel parses the configured logging level into one understood by
// logrus.
func (c *Config) setLogLevel() {
	level, err := logrus.ParseLevel(c.UnparsedLogLevel)

	if err != nil {
		level = logrus.InfoLevel
	}

	c.LogLevel = level
}
