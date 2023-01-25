package listeners

import (
	"fmt"
	"strconv"
	"strings"

	t "github.com/ryanlake6/squash-slack-bot/types"
	"github.com/shomali11/slacker"
)

func (c *Client) FindNextMatch() {
	c.Bot.Command("next match <playerName>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			playerName := strings.ToLower(request.Param("playerName"))

			// gets ordered recent matches
			pastmatches, errFromGetRecentMatchesByDate := c.getRecentMatchesByDate(playerName)

			// get players current position
			playerRanking, errFromGetPlayerRanking := c.getPlayerRanking(playerName)

			var resp string
			// Checking if first two sql calls are valid/work
			if errFromGetPlayerRanking == nil && errFromGetRecentMatchesByDate == nil {
				if pastmatches[0].Player1 == playerName && pastmatches[0].Player1PrevPos > pastmatches[0].Player2PrevPos ||
					pastmatches[0].Player2 == playerName && pastmatches[0].Player1PrevPos < pastmatches[0].Player2PrevPos {
						playerToPlayNext, err := c.getPlayerBasedOnRanking(int(playerRanking.Position) + 1)
						if (err != nil) {
							resp = err.Error()
						} else {
							resp = "You should play DOWN next against "  + playerToPlayNext.FirstName
						}
				} else {
					playerToPlayNext, err := c.getPlayerBasedOnRanking(int(playerRanking.Position) - 1)
					if (err != nil) {
						resp = err.Error()
					} else {
						resp = "You should play UP next against " + playerToPlayNext.FirstName
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

// gets the ordered most recent matches by date of the given player
func (c *Client) getRecentMatchesByDate(playerName string) ([]t.PastMatch, error) {
	rows, _ := c.Database.Query("SELECT * from pastmatches WHERE player1='" + playerName + "' OR " + "player2='" + playerName + "' ORDER BY date DESC")
		pastmatches := []t.PastMatch{}
		defer rows.Close()
		for rows.Next() {
			var pm t.PastMatch
			rows.Scan(&pm.Player1, &pm.Player2, &pm.Winner, &pm.Player1PrevPos, &pm.Player2PrevPos, &pm.Date)
			pastmatches = append(pastmatches, pm)
		}
		if (len(pastmatches) == 0) {
			return nil, fmt.Errorf("No matches have been played by this player")
		}

	return pastmatches, nil
}


func (c *Client) getPlayerRanking(playerName string) (t.Player, error) {
	rows, _ := c.Database.Query("SELECT * FROM rankings WHERE firstName='" + playerName + "'")
	players := []t.Player{}
	defer rows.Close()
	for rows.Next() {
		var p t.Player
		rows.Scan(&p.Position, &p.FirstName)
		players = append(players, p)
	}

	if (len(players) == 0) {
		return *&t.Player{}, fmt.Errorf("Can't get ranking if there player doesn't exist in database")
	}
	return players[0], nil
}

func (c *Client) getPlayerBasedOnRanking(rankingNumber int) (t.Player, error) {
	// Checking if player is at top or bottom of ladder
	ladder := c.retreiveLadderFromDatabase()
	if (rankingNumber == 0 || len(ladder) < rankingNumber) {
		return *&t.Player{}, fmt.Errorf("You are top or bottom of the ladder and thus can only play one direction")
	}

	rows, _ := c.Database.Query("SELECT * FROM rankings WHERE position=" + strconv.Itoa(rankingNumber))
	players := []t.Player{}
	defer rows.Close()
	for rows.Next() {
		var p t.Player
		rows.Scan(&p.Position, &p.FirstName)
		players = append(players, p)
	}
	if (len(players) == 0) {
		return *&t.Player{}, fmt.Errorf("Can't get this player based on rank given") 
	}
	return players[0], nil
}