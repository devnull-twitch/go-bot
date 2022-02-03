package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/devnull-twitch/go-bot/internal/commands"
	"github.com/devnull-twitch/go-bot/internal/modules"
	"github.com/devnull-twitch/go-tmi"
	"github.com/google/go-github/v42/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	godotenv.Load(".env.yaml")

	bot, err := tmi.New(os.Getenv("USERNAME"), os.Getenv("TOKEN"), os.Getenv("CHANNEL"), os.Getenv("COMMAND_MARK"))
	if err != nil {
		log.Fatal(err)
	}

	var httpClient *http.Client
	if os.Getenv("GH_PAT") != "" {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GH_PAT")},
		)
		httpClient = oauth2.NewClient(ctx, ts)
	}
	ghClient := github.NewClient(httpClient)

	bot.AddCommand(commands.RandChatter())
	bot.AddCommand(commands.GithubDataCommand(ghClient))
	bot.AddCommand(commands.GithubMakeTagCommand(ghClient))
	bot.AddCommand(commands.ListCommandsCommand())
	for _, c := range commands.TextResponses() {
		bot.AddCommand(c)
	}

	bot.AddModule(modules.TimedMessageMod([]string{"This", "is", "a", "test"}, os.Getenv("CHANNEL")))

	if err := bot.Run(); err != nil {
		log.Fatal(err)
	}
}
