package listeners

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
)

func (c *Client) InputMatch() {
	c.Bot.Command("match <player1> beats <player2>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			player1 := request.Param("player1")
			player2 := request.Param("player2")

			rows, _ := c.Database.Query("SELECT * FROM rankings WHERE firstName='" + player1 + "'")
			players1 := []Player{}
			defer rows.Close()
			for rows.Next() {
				var p Player
				rows.Scan(&p.Position, &p.FirstName)
				players1 = append(players1, p)
			}

			rows, _ = c.Database.Query("SELECT * FROM rankings WHERE firstName='" + player2 + "'")
			defer rows.Close()
			players2 := []Player{}
			for rows.Next() {
				var p Player
				rows.Scan(&p.Position, &p.FirstName)
				players2 = append(players2, p)
			}
			player1Position := players1[0].Position
			player2Position := players2[0].Position
			player1Name := players1[0].FirstName
			player2Name := players2[0].FirstName
			// response.Reply(strconv.Itoa(int(player1Position)) + strconv.Itoa(int(player2Position)) + player1Name + player2Name)
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