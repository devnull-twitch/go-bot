package commands

import (
	"fmt"

	"github.com/devnull-twitch/go-tmi/pkg/tmi"
)

func cmdLine(cmd tmi.Command) string {
	params := " "
	for _, p := range cmd.Params {
		if p.Required {
			params += fmt.Sprintf("<%s> ", p.Name)
		} else {
			params += fmt.Sprintf("[<%s>] ", p.Name)
		}
	}
	if cmd.AllowRestParams {
		params += "[...] "
	}

	return fmt.Sprintf("!%s%s- %s | ", cmd.Name, params, cmd.Description)
}

func ListCommandsCommand() tmi.Command {
	return tmi.Command{
		Name: "commands",
		Handler: func(client *tmi.Client, args tmi.CommandArgs) *tmi.OutgoingMessage {
			msg := ""
			for _, cmd := range client.ListCommands() {
				if !cmd.RequiresBroadcasterOrMod {
					msg += cmdLine(cmd)
				}
			}

			if args.UserIsBroadcasterOrMod {
				msg += "Mod commands ||| "
				for _, cmd := range client.ListCommands() {
					if cmd.RequiresBroadcasterOrMod {
						msg += cmdLine(cmd)
					}
				}
			}

			return &tmi.OutgoingMessage{
				Message: msg[0 : len(msg)-3],
			}
		},
	}
}
