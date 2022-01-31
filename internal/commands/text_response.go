package commands

import (
	"io/ioutil"

	"github.com/devnull-twitch/go-tmi/pkg/tmi"
	"gopkg.in/yaml.v2"
)

type (
	CommandConfig struct {
		Command  string `yaml:"command"`
		Response string `yaml:"response"`
	}
	Responses struct {
		Responses []CommandConfig `yaml:"responses,flow"`
	}
)

func TextResponses() []tmi.Command {
	config := &Responses{}
	inputBytes, err := ioutil.ReadFile("texts.yaml")
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(inputBytes, config)

	cmds := make([]tmi.Command, len(config.Responses))
	for i, c := range config.Responses {
		respCopy := c.Response
		cmds[i] = tmi.Command{
			Name: c.Command,
			Handler: func(client *tmi.Client, args tmi.CommandArgs) *tmi.OutgoingMessage {
				return &tmi.OutgoingMessage{
					Message: respCopy,
				}
			},
		}
	}

	return cmds
}
