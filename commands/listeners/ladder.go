package listeners

import (
	"strconv"

	sqlQueries "github.com/ryanlake6/squash-slack-bot/database"
	"github.com/shomali11/slacker"
)

type Client struct {
	Bot *slacker.Slacker
	Database *sqlQueries.Database
}

func (c *Client) GetLadder() {
	c.Bot.Command("ladder", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			players := c.Database.RetreiveLadderFromDatabase()
			
			temp := ""
			for i := 0; i < len(players); i++ {
				temp = temp + "#" + strconv.FormatInt(players[i].Position, 10) + " " + players[i].FirstName + "\n"
			}
			response.Reply(temp)
		},
	})
}
