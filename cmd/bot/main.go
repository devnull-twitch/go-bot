package main

import (
	"log"
	"os"

	"github.com/devnull-twitch/go-tmi/internal/commands"
	"github.com/devnull-twitch/go-tmi/pkg/tmi"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.yaml")

	bot, err := tmi.New(os.Getenv("USERNAME"), os.Getenv("TOKEN"), os.Getenv("CHANNEL"), os.Getenv("COMMAND_MARK"))
	if err != nil {
		log.Fatal(err)
	}

	bot.AddCommand(commands.RandChatter())
	bot.AddCommand(commands.WaitResponse())
	for _, c := range commands.TextResponses() {
		bot.AddCommand(c)
	}

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
