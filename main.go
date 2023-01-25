package main

import (
	"context"
	"log"
	"os"

	"github.com/ryanlake6/squash-slack-bot/commands"
	"github.com/ryanlake6/squash-slack-bot/commands/listeners"
	"github.com/ryanlake6/squash-slack-bot/database"
	"github.com/shomali11/slacker"
)

func main() {
	// Connect to sql database
	db := database.DatabaseConnect()

	// Creating the bot instance based on the Squash slack info/tokens
	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	// Creating all the bot commands
	database := &database.Database{Database: db}
	listeners := &listeners.Client{Bot: bot, Database: database}
	c := &commands.Client{CommandEvents: listeners}
	c.CreateCommands()

	// When program gets cancelled at terminal or any other reason (shuts down the bot)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Kills the bot if an error occurs on a call
	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}


