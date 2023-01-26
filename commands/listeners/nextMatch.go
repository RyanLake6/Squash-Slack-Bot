package listeners

import (
	"strings"

	"github.com/shomali11/slacker"
)

func (c *Client) FindNextMatch() {
	c.Bot.Command("next match <playerName>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			playerName := strings.ToLower(request.Param("playerName"))

			// gets ordered recent matches
			pastmatches, errFromGetRecentMatchesByDate := c.Database.GetRecentMatchesByDate(playerName)

			// get players current position
			playerRanking, errFromGetPlayerRanking := c.Database.GetPlayerRanking(playerName)

			var resp string
			// Checking if first two sql calls are valid/work
			if errFromGetPlayerRanking == nil && errFromGetRecentMatchesByDate == nil {
				if pastmatches[0].Player1 == playerName && pastmatches[0].Player1PrevPos > pastmatches[0].Player2PrevPos ||
					pastmatches[0].Player2 == playerName && pastmatches[0].Player1PrevPos < pastmatches[0].Player2PrevPos {
						playerToPlayNext, edgeCaseRanking, err := c.Database.GetPlayerBasedOnRanking(int(playerRanking.Position) + 1)
						if (err != nil) {
							resp = err.Error()
						} else if edgeCaseRanking != "" { // For the case of the number 1 postion not able to challenge up
							playerToPlayNext, _, _ = c.Database.GetPlayerBasedOnRanking(int(playerRanking.Position) - 1)
							resp = edgeCaseRanking + playerToPlayNext.FirstName
						} else {
							resp = "You should play DOWN next against "  + playerToPlayNext.FirstName
						}
				} else {
					playerToPlayNext, edgeCaseRanking, err := c.Database.GetPlayerBasedOnRanking(int(playerRanking.Position) - 1)
					if (err != nil) {
						resp = err.Error()
					} else if edgeCaseRanking != "" { // For the case of the lowest player not able to challenge down
						playerToPlayNext, _, _ = c.Database.GetPlayerBasedOnRanking(int(playerRanking.Position) + 1)
						resp = edgeCaseRanking + playerToPlayNext.FirstName
					} else {
						resp = "You should play UP next against " + playerToPlayNext.FirstName + edgeCaseRanking
					}
				}
			}

			if (len(resp) == 0) {
				response.Reply("I don't know you haven't inputted a match yet")
			} else {
				response.Reply(resp)
			}
		},
	})
}