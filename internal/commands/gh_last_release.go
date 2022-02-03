package commands

import (
	"context"
	"fmt"

	"github.com/devnull-twitch/go-bot/pkg/tmi"
	"github.com/google/go-github/v42/github"
	"github.com/sirupsen/logrus"
)

func GithubDataCommand(ghClient *github.Client) tmi.Command {
	return tmi.Command{
		Name: "release",
		Handler: func(client *tmi.Client, args tmi.CommandArgs) *tmi.OutgoingMessage {
			rel, resp, err := ghClient.Repositories.GetLatestRelease(context.Background(), "devnull-twitch", "2donlinerpg")
			if err != nil {
				logrus.WithError(err).Error("unable to load latest release")
				return nil
			}

			if resp.StatusCode != 200 {
				logrus.WithFields(logrus.Fields{
					"response": resp.Status,
				}).Error("unable to load latest release")
				return nil
			}

			msg := ""
			for _, r := range rel.Assets {
				msg += fmt.Sprintf("%s => %s | ", *r.Name, *r.BrowserDownloadURL)
			}

			if len(msg) > 3 {
				msg = msg[:len(msg)-3]
				return &tmi.OutgoingMessage{
					Message: msg,
				}
			}
			return nil
		},
	}
}
