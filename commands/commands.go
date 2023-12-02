package commands

import (
	commandfunctions "github.com/ryanlake6/squash-slack-bot/commands/listeners"
)

type Client struct {
	CommandEvents *commandfunctions.Client
}

// Main function to create all the bot.commands 
func (c *Client) CreateCommands() {
	// display incoming commands within terminal (debugging)
	go c.CommandEvents.PrintCommandEvents(c.CommandEvents.Bot.CommandEvents())

	// outputs the ladder
	c.CommandEvents.GetLadder()

	// allows to input match into system
	c.CommandEvents.InputMatch()

	// Gets last matches of the player
	c.CommandEvents.GetLastMatches()

	// Find next match player should play
	c.CommandEvents.FindNextMatch()

	// Manually alter the ladder (admin access only)
	c.CommandEvents.AlterLadder()
}