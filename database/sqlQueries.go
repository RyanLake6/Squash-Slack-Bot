package database

import (
	"database/sql"
	"fmt"
	"strconv"

	t "github.com/ryanlake6/squash-slack-bot/types"
)

type Database struct {
	Database *sql.DB
}

func (d *Database) RetreiveLadderFromDatabase() []t.Player {
	rows, _ := d.Database.Query("SELECT * from rankings ORDER BY position ASC")
	players := []t.Player{}
	defer rows.Close()
	for rows.Next() {
		var p t.Player
		rows.Scan(&p.Position, &p.FirstName)
		players = append(players, p)
	}

	return players
}

// gets the ordered most recent matches by date of the given player
func (d *Database) GetRecentMatchesByDate(playerName string) ([]t.PastMatch, error) {
	rows, _ := d.Database.Query("SELECT * from pastmatches WHERE player1='" + playerName + "' OR " + "player2='" + playerName + "' ORDER BY date DESC")
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

// Get a single players ranking
func (d *Database) GetPlayerRanking(playerName string) (t.Player, error) {
	rows, _ := d.Database.Query("SELECT * FROM rankings WHERE firstName='" + playerName + "'")
	players := []t.Player{}
	defer rows.Close()
	for rows.Next() {
		var p t.Player
		rows.Scan(&p.Position, &p.FirstName)
		players = append(players, p)
	}

	if (len(players) == 0) {
		return *&t.Player{}, fmt.Errorf("Can't get ranking if the player doesn't exist in database")
	}
	return players[0], nil
}

// Returns the player type based on the rank given
func (d *Database) GetPlayerBasedOnRanking(rankingNumber int) (t.Player, string, error) {
	// Checking if player is at top or bottom of ladder
	ladder := d.RetreiveLadderFromDatabase()
	edgeCaseRanking := ""
	if (rankingNumber == 0) {
		rankingNumber += 2
		edgeCaseRanking = "Since you are number 1 you always will always play DOWN against "
	} else if (rankingNumber > len(ladder)) {
		rankingNumber -= 2
		edgeCaseRanking = "Since you are at the bottom of the ladder you always will play UP against " 
	}

	rows, _ := d.Database.Query("SELECT * FROM rankings WHERE position=" + strconv.Itoa(rankingNumber))
	players := []t.Player{}
	defer rows.Close()
	for rows.Next() {
		var p t.Player
		rows.Scan(&p.Position, &p.FirstName)
		players = append(players, p)
	}
	if (len(players) == 0) {
		return *&t.Player{}, "", fmt.Errorf("Can't get this player based on rank given") 
	}
	return players[0], edgeCaseRanking, nil
}


// updates the ladder after a match has been completed
func (d *Database) UpdateLadder(winner string, winnersOldPosition int, loser string, losersOldPosition int) {
	_, err := d.Database.Exec("UPDATE rankings SET position='" + strconv.Itoa(winnersOldPosition - 1) + "'" + "WHERE firstname='" + winner + "'")
	if err != nil {
		panic(err)
	}
	_, err = d.Database.Exec("UPDATE rankings SET position='" + strconv.Itoa(losersOldPosition + 1) + "'" + "WHERE firstname='" + loser + "'")
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated ladder!")
}

// records the match in the pastmatches table
func (d *Database) RecordMatch(winner int, player1 string, player2 string, player1PrevPos int, player2PrevPos int) {
	player1Temp := "'" + player1 + "', "
	player2Temp := "'" + player2 + "', "
	winnerTemp := "'" + strconv.Itoa(winner) + "', "
	player1PrevPosTemp := "'" + strconv.Itoa(player1PrevPos) + "', "
	player2PrevPosTemp := "'" + strconv.Itoa(player2PrevPos) + "', "
	// currentTime := "'" + time.Now().Format(time.RFC3339) + "'"
	_, err := d.Database.Exec("INSERT INTO pastmatches VALUES (" + player1Temp + player2Temp + winnerTemp + player1PrevPosTemp + player2PrevPosTemp + "'2023-01-21'" + ")")
	if err != nil {
		panic(err)
	}
}

// Gets all the past matches of the player given
func (d *Database) GetPastMatches(playerName string) []t.PastMatch {
	rows, _ := d.Database.Query("SELECT * from pastmatches WHERE player1='" + playerName + "' OR " + "player2='" + playerName + "'")
	pastmatches := []t.PastMatch{}
	defer rows.Close()
	for rows.Next() {
		var pm t.PastMatch
		rows.Scan(&pm.Player1, &pm.Player2, &pm.Winner, &pm.Player1PrevPos, &pm.Player2PrevPos, &pm.Date)
		pastmatches = append(pastmatches, pm)
	}

	return pastmatches
}