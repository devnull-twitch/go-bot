package commands

import (
	"time"

	"github.com/devnull-twitch/go-tmi/pkg/tmi"
)

func WaitResponse() tmi.Command {
	return tmi.Command{
		Name: "late",
		Handler: func(client *tmi.Client, args tmi.CommandArgs) *tmi.OutgoingMessage {
			go func() {
				time.Sleep(5 * time.Second)
				client.Send(&tmi.OutgoingMessage{
					Message:     "Sorry .. wasnt listening",
					Channel:     args.Channel,
					ParentID:    args.ParentID,
					SendAsReply: true,
				})
			}()

			return nil
		},
	}
}
