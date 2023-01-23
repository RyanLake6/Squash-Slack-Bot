package main

import (
	"context"
	"log"
	"os"

	"github.com/ryanlake6/squash-slack-bot/commands"
	"github.com/ryanlake6/squash-slack-bot/database"
	"github.com/shomali11/slacker"
)

func main() {
	// Connect to sql database
	db := database.DatabaseConnect()

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	// Creating all the bot commands
	c := commands.Client{Database: db, Bot: bot}
	c.CreateCommands()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}


