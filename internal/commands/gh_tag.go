package commands

import (
	"context"
	"fmt"
	"regexp"

	"github.com/devnull-twitch/go-bot/pkg/tmi"
	"github.com/google/go-github/v42/github"
	"github.com/sirupsen/logrus"
)

func strPtr(s string) *string {
	return &s
}

func GithubMakeTagCommand(ghClient *github.Client) tmi.Command {
	return tmi.Command{
		Name: "tag",
		Params: []tmi.Parameter{
			{
				Name:     "name",
				Required: true,
				Validate: func(s string) bool {
					matched, _ := regexp.MatchString(`^v(\d+)\.(\d+).(\d+)$`, s)
					return matched
				},
			},
		},
		RequiresBroadcasterOrMod: true,
		AllowRestParams:          false,
		Handler: func(client *tmi.Client, args tmi.CommandArgs) *tmi.OutgoingMessage {
			ghBranch, branchResp, err := ghClient.Repositories.GetBranch(context.Background(), "devnull-twitch", "2donlinerpg", "master", false)
			if err != nil || ghBranch == nil || ghBranch.Commit == nil {
				logrus.WithFields(logrus.Fields{
					"err":    err,
					"status": branchResp.Status,
				}).Error("unable to fetch branch")
				return &tmi.OutgoingMessage{
					Message: "Error reading latest commit NotLikeThis",
				}
			}

			tagRef := fmt.Sprintf("refs/tags/%s", args.Parameters["name"])
			ghRef, refResp, err := ghClient.Git.CreateRef(context.Background(), "devnull-twitch", "2donlinerpg", &github.Reference{
				Ref: &tagRef,
				Object: &github.GitObject{
					Type: strPtr("commit"),
					SHA:  ghBranch.Commit.SHA,
					URL:  ghBranch.Commit.URL,
				},
			})
			if err != nil || ghBranch == nil || ghBranch.Commit == nil {
				logrus.WithFields(logrus.Fields{
					"err":    err,
					"status": refResp.Status,
				}).Error("unable to create tag reference")
				return &tmi.OutgoingMessage{
					Message: "Error creating reference NotLikeThis",
				}
			}

			ghTag, tagResp, err := ghClient.Git.CreateTag(context.Background(), "devnull-twitch", "2donlinerpg", &github.Tag{
				Tag:     strPtr(args.Parameters["name"]),
				Message: strPtr(args.Parameters["name"]),
				Object:  ghRef.GetObject(),
			})
			if err != nil || *ghTag.Tag == "" {
				logrus.WithFields(logrus.Fields{
					"err":    err,
					"status": tagResp.Status,
				}).Error("unable to create tag")
				return &tmi.OutgoingMessage{
					Message: "Error creating tag NotLikeThis",
				}
			}

			return &tmi.OutgoingMessage{
				Message: fmt.Sprintf("Done! Commit %s is now a tag named %s!", *ghBranch.Commit.SHA, *ghTag.Tag),
			}
		},
	}
}
