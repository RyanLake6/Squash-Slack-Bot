package listeners

import (
	"database/sql"
	"strconv"

	t "github.com/ryanlake6/squash-slack-bot/types"
	"github.com/shomali11/slacker"
)

type Client struct {
	Database *sql.DB
	Bot *slacker.Slacker
}

func (c *Client) GetLadder() {
	c.Bot.Command("ladder", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			rows, _ := c.Database.Query("SELECT * from rankings ORDER BY position ASC")
			players := []t.Player{}
			defer rows.Close()
			for rows.Next() {
				var p t.Player
				rows.Scan(&p.Position, &p.FirstName)
				players = append(players, p)
			}
			temp := ""
			for i := 0; i < len(players); i++ {
				temp = temp + "#" + strconv.FormatInt(players[i].Position, 10) + " " + players[i].FirstName + "\n"
			}
			response.Reply(temp)
		},
	})
}