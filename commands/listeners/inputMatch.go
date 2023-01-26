package listeners

import (
	"github.com/shomali11/slacker"
)

func (c *Client) InputMatch() {
	c.Bot.Command("match <player1> beats <player2>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			player1 := request.Param("player1")
			player2 := request.Param("player2")
			var resp string

			player1Ranking, err := c.Database.GetPlayerRanking(player1)
			if (err != nil) {
				resp = err.Error()
			}
			player2Ranking, err := c.Database.GetPlayerRanking(player2)
			if (err != nil) {
				resp = err.Error()
			}

			player1Position := player1Ranking.Position
			player2Position := player2Ranking.Position
			player1Name := player1Ranking.FirstName
			player2Name := player2Ranking.FirstName
			if (resp == "") {
				if int(player1Position) == int(player2Position)-1 && int(player1Position) < int(player2Position) {
					resp = "Great job " + player1Name + "! \nThe ladder won't change and this match has been recorded"
					c.Database.RecordMatch(1, player1Name, player2Name, int(player1Position), int(player2Position))
				} else if int(player1Position) == int(player2Position)+1 && player1Position > player2Position {
					resp = "Great job " + player1Name + "! \nThe ladder will be affected and this match has been recorded"
					c.Database.UpdateLadder(player1Name, int(player1Position), player2Name, int(player2Position))
					c.Database.RecordMatch(2, player1Name, player2Name, int(player1Position), int(player2Position))
				} else {
					response.Reply("uh oh I can't support players beating eachother not next to eachother in the ladder :(")
				}
			}
			response.Reply(resp)
		},
	})
}
