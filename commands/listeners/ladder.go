package listeners

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
)

type Client struct {
	Database *sql.DB
	Bot *slacker.Slacker
}

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

func (c *Client) GetLadder() {
	c.Bot.Command("ladder", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			rows, _ := c.Database.Query("SELECT * from rankings ORDER BY position ASC")
			players := []Player{}
			defer rows.Close()
			for rows.Next() {
				var p Player
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