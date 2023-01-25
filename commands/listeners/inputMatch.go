package listeners

import (
	"fmt"
	"strconv"
	"time"

	t "github.com/ryanlake6/squash-slack-bot/types"
	"github.com/shomali11/slacker"
)

func (c *Client) InputMatch() {
	c.Bot.Command("match <player1> beats <player2>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			player1 := request.Param("player1")
			player2 := request.Param("player2")

			player1Ranking, player2Ranking := c.getPlayersRankings(player1, player2)

			player1Position := player1Ranking.Position
			player2Position := player2Ranking.Position
			player1Name := player1Ranking.FirstName
			player2Name := player2Ranking.FirstName
			if int(player1Position) == int(player2Position)-1 && int(player1Position) < int(player2Position) {
				response.Reply("Great job " + player1Name + "! \nThe ladder won't change and this match has been recorded")
				c.recordMatch(1, player1Name, player2Name, int(player1Position), int(player2Position))
			} else if int(player1Position) == int(player2Position)+1 && player1Position > player2Position {
				response.Reply("Great job " + player1Name + "! \nThe ladder will be affected and this match has been recorded")
				c.updateLadder(player1Name, int(player1Position), player2Name, int(player2Position))
				c.recordMatch(2, player1Name, player2Name, int(player1Position), int(player2Position))
			} else {
				response.Reply("uh oh I can't support players beating eachother not next to eachother in the ladder :(")
			}
		},
	})
}


func (c *Client) getPlayersRankings(player1 string, player2 string) (t.Player, t.Player) {
	rows, _ := c.Database.Query("SELECT * FROM rankings WHERE firstName='" + player1 + "'")
	players1 := []t.Player{}
	defer rows.Close()
	for rows.Next() {
		var p t.Player
		rows.Scan(&p.Position, &p.FirstName)
		players1 = append(players1, p)
	}

	rows, _ = c.Database.Query("SELECT * FROM rankings WHERE firstName='" + player2 + "'")
	defer rows.Close()
	players2 := []t.Player{}
	for rows.Next() {
		var p t.Player
		rows.Scan(&p.Position, &p.FirstName)
		players2 = append(players2, p)
	}

	return players1[0], players2[0]
}

// records the match in the pastmatches table
func (c *Client) recordMatch(winner int, player1 string, player2 string, player1PrevPos int, player2PrevPos int) {
	player1Temp := "'" + player1 + "', "
	player2Temp := "'" + player2 + "', "
	winnerTemp := "'" + strconv.Itoa(winner) + "', "
	player1PrevPosTemp := "'" + strconv.Itoa(player1PrevPos) + "', "
	player2PrevPosTemp := "'" + strconv.Itoa(player2PrevPos) + "', "
	currentTime := "'" + time.Now().Format(time.RFC3339) + "'"
	_, err := c.Database.Exec("INSERT INTO pastmatches VALUES (" + player1Temp + player2Temp + winnerTemp + player1PrevPosTemp + player2PrevPosTemp + currentTime + ")")
	if err != nil {
		panic(err)
	}
}

// updates the ladder after a match has been completed
func (c *Client) updateLadder(winner string, winnersOldPosition int, loser string, losersOldPosition int) {
	_, err := c.Database.Exec("UPDATE rankings SET position='" + strconv.Itoa(winnersOldPosition - 1) + "'" + "WHERE firstname='" + winner + "'")
	if err != nil {
		panic(err)
	}
	_, err = c.Database.Exec("UPDATE rankings SET position='" + strconv.Itoa(losersOldPosition + 1) + "'" + "WHERE firstname='" + loser + "'")
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated ladder!")
}