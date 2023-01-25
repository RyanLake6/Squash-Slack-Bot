package listeners

import (
	"strconv"
	"strings"

	"github.com/shomali11/slacker"
)

func (c *Client) GetLastMatches() {
	c.Bot.Command("last matches <playerName>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			playerName := strings.ToLower(request.Param("playerName"))

			pastmatches := c.Database.GetPastMatches(playerName)		

			// building the response
			temp := ""
			for i := 0; i < len(pastmatches); i++ {
				// Finding who is the other player
				var otherPlayer string
				if playerName == pastmatches[i].Player1 {
					otherPlayer = pastmatches[i].Player2
				} else if playerName == pastmatches[i].Player2 {
					otherPlayer = pastmatches[i].Player1
				}
				// Finding if the player who was called on won or lost
				var whoWon string
				if pastmatches[i].Winner == 1 && pastmatches[i].Player1 == playerName {
					whoWon = "WON"
				} else {
					whoWon = "LOST"
				}
				// Finding if the player who was called on played up or down
				var upOrDown string
				if pastmatches[i].Player1 == playerName && pastmatches[i].Player1PrevPos > pastmatches[i].Player2PrevPos ||
					pastmatches[i].Player2 == playerName && pastmatches[i].Player1PrevPos < pastmatches[i].Player2PrevPos {
					upOrDown = "UP"
				} else {
					upOrDown = "DOWN"
				}

				date := pastmatches[i].Date.Month().String() + "-" + strconv.Itoa(pastmatches[i].Date.Day()) + "-" + strconv.Itoa(pastmatches[i].Date.Year())
				// Building the final output
				temp = temp + "#" + strconv.Itoa(i+1) + " Against " + otherPlayer + " you played " + upOrDown + " and " + whoWon + " on " + date + "\n"
			}

			// Sending reponse if there is/isnt matches associated with the player
			if len(temp) == 0 {
				response.Reply("Sorry, you don't have any stored matches")
			} else {
				response.Reply(temp)
			}
		},
	})
}


