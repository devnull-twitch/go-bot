package commands

import (
	"fmt"
	"math/rand"

	"github.com/devnull-twitch/go-bot/pkg/tmi"
)

func RandChatter() tmi.Command {
	return tmi.Command{
		Name:                     "rnd_chatter",
		Description:              "Picks a random chatter",
		RequiresBroadcasterOrMod: true,
		Handler: func(client *tmi.Client, args tmi.CommandArgs) *tmi.OutgoingMessage {
			all := client.Chatters()

			if len(all) <= 1 {
				return &tmi.OutgoingMessage{Message: "That's boring"}
			}

			target := rand.Intn(len(all) - 1)
			return &tmi.OutgoingMessage{
				Message: fmt.Sprintf("Picked %s", all[target]),
			}
		},
	}
}
