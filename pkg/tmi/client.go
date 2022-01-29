package tmi

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/irc.v3"
)

type Client struct {
	ircClient *irc.Client
	userList  map[string]int64
	commands  map[string]Command
}

type IncomingCommand struct {
	Broadcaster bool
	Mod         bool
	User        string
	Channel     string
	Command     string
	Params      []string
	MsgID       string
}

type OutgoingMessage struct {
	Message     string
	Channel     string
	ParentID    string
	SendAsReply bool
}

func New(username, token, channel, commandMarkerChar string) (*Client, error) {
	conn, err := tls.Dial("tcp", "irc.chat.twitch.tv:6697", &tls.Config{})
	if err != nil {
		return nil, err
	}

	tmiClient := &Client{
		userList: make(map[string]int64),
		commands: make(map[string]Command),
	}

	ircConfig := irc.ClientConfig{
		Nick: username,
		Pass: token,
		User: username,
		Name: username,
		Handler: irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
			if m.Command == "JOIN" && m.Prefix != nil && m.Prefix.User != "" {
				tmiClient.handleJoinFrom(m.Prefix.User)
				return
			}
			if m.Command == "PART" && m.Prefix != nil && m.Prefix.User != "" {
				tmiClient.handleUserPart(m.Prefix.User)
				return
			}
			if m.Command == "375" || m.Command == "376" || m.Command == "372" {
				// on MODT end we make join request
				if m.Command == "376" {
					c.Write(fmt.Sprintf("JOIN #%s", channel))
				}
				// MODT commands
				return
			}
			if m.Command == "001" || m.Command == "002" || m.Command == "003" || m.Command == "004" {
				// Welcome spam
				return
			}

			if m.Command == "PRIVMSG" {
				if m.Param(1)[0:1] != commandMarkerChar {
					return
				}

				incoming := &IncomingCommand{
					Channel: m.Param(0),
				}

				params := strings.Split(m.Param(1)[1:], " ")
				incoming.Command = params[0]
				incoming.Params = params[1:]

				if len(m.Tags) > 0 {
					for tagName, tagValue := range m.Tags {
						if tagName == "badges" && strings.Contains(string(tagValue), "broadcaster/1") {
							incoming.Broadcaster = true
						}
						if tagName == "mod" && tagValue == "1" {
							incoming.Mod = true
						}
						if tagName == "id" {
							incoming.MsgID = string(tagValue)
						}
					}
				}

				tmiClient.handleCommand(incoming)
			}
		}),
	}

	ircClient := irc.NewClient(conn, ircConfig)
	ircClient.CapRequest("twitch.tv/membership", true)
	ircClient.CapRequest("twitch.tv/tags", true)

	tmiClient.ircClient = ircClient

	return tmiClient, nil
}

func (c *Client) Run() error {
	logrus.Info("Starting bot")
	return c.ircClient.Run()
}
