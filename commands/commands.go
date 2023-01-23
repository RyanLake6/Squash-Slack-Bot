package commands

import (
	"database/sql"
	"time"

	"github.com/shomali11/slacker"
)

type Player struct {
	Position  int64
	FirstName string
}

type PastMatch struct {
	Player1 string
	Player2 string
	Winner int
	Player1PrevPos int
	Player2PrevPos int
	Date time.Time
}

type Client struct {
	Database *sql.DB
	Bot *slacker.Slacker
}

// Main function to create all the bot.commands 
func (c *Client) CreateCommands() {
	// outputs the ladder
	c.GetLadder()

	// allows to input match into system
	c.inputMatch()

	c.getLastMatches()


}