package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
	}
}

var db *sql.DB

func main() {
		godotenv.Load()
	    // Capture connection properties.
		cfg := mysql.Config{
			ParseTime: true,
			User:   os.Getenv("DBUSER"),
			Passwd: os.Getenv("DBPASS"),
			Net:    "tcp",
			Addr:   "127.0.0.1:3306",
			DBName: "ladder",
		}
		// Get a database handle.
		var err error
		db, err = sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			log.Fatal(err)
		}
	
		pingErr := db.Ping()
		if pingErr != nil {
			log.Fatal(pingErr)
		}
		fmt.Println("Connected to sql!")


	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())


	type Player struct {
		position int64
		firstName string
	}
	bot.Command("ladder", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			rows, _ := db.Query("SELECT * from rankings ORDER BY position ASC")
			players := []Player{}
			defer rows.Close()
			for rows.Next() {
				var p Player
				rows.Scan(&p.position, &p.firstName)
				players = append(players, p)
			}
			temp := ""
			for i := 0; i < len(players); i++ {
				temp = temp + "#" + strconv.FormatInt(players[i].position, 10) + " " + players[i].firstName + "\n"
			}
			response.Reply(temp)
		},
	})

	bot.Command("match <player1> beats <player2>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			player1 := request.Param("player1")
			player2 := request.Param("player2")

			rows, _ := db.Query("SELECT * FROM rankings WHERE firstName='" + player1 + "'")
			players1 := []Player{}
			defer rows.Close()
			for rows.Next() {
				var p Player
				rows.Scan(&p.position, &p.firstName)
				players1 = append(players1, p)
			}

			rows, _ = db.Query("SELECT * FROM rankings WHERE firstName='" + player2 + "'")
			defer rows.Close()
			players2 := []Player{}
			for rows.Next() {
				var p Player
				rows.Scan(&p.position, &p.firstName)
				players2 = append(players2, p)
			}
			player1Position := players1[0].position
			player2Position := players2[0].position
			player1Name := players1[0].firstName
			player2Name := players2[0].firstName
			// response.Reply(strconv.Itoa(int(player1Position)) + strconv.Itoa(int(player2Position)) + player1Name + player2Name)
			if (int(player1Position) == int(player2Position) - 1 && int(player1Position) < int(player2Position)) {
				response.Reply("Great job " + player1Name + "! \nThe ladder won't change and this match has been recorded")
				recordMatch(1, player1Name, player2Name, int(player1Position), int(player2Position))
			} else if (int(player1Position) == int(player2Position) + 1 && player1Position > player2Position) {
				response.Reply("Great job " + player1Name + "! \nThe ladder will be affected and this match has been recorded")
				updateLadder(player1Name, int(player1Position), player2Name, int(player2Position))
				recordMatch(2, player1Name, player2Name, int(player1Position), int(player2Position))
			} else {
				response.Reply("uh oh I can't support players beating eachother not next to eachother in the ladder :(")
			}
		},
	})

	type PastMatch struct {
		player1 string
		player2 string
		winner int
		player1PrevPos int
		player2PrevPos int
		date time.Time
	}
	bot.Command("last matches <playerName>", &slacker.CommandDefinition{
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			playerName := request.Param("playerName")

			rows, _ := db.Query("SELECT * from pastmatches WHERE player1='" + playerName + "' OR " + "player2='" + playerName + "'")
			pastmatches := []PastMatch{}
			defer rows.Close()
			for rows.Next() {
				var pm PastMatch
				rows.Scan(&pm.player1, &pm.player2, &pm.winner, &pm.player1PrevPos, &pm.player2PrevPos, &pm.date)
				pastmatches = append(pastmatches, pm)
			}

			// building the response
			temp := ""
			for i := 0; i < len(pastmatches); i++ {
				// Finding who is the other player
				var otherPlayer string
				if (playerName == pastmatches[i].player1) {
					otherPlayer = pastmatches[i].player2
				} else {
					otherPlayer = pastmatches[i].player1
				}
				// Finding if the player who was called on won or lost
				var whoWon string
				if (pastmatches[i].winner == 1 && pastmatches[i].player1 == playerName) {
					whoWon = "WON"
				} else {
					whoWon = "LOST"
				}
				// Finding if the player who was called on played up or down
				var upOrDown string
				if (pastmatches[i].player1 == playerName && pastmatches[i].player1PrevPos > pastmatches[i].player2PrevPos ||
					pastmatches[i].player2 == playerName && pastmatches[i].player1PrevPos < pastmatches[i].player2PrevPos) {
					upOrDown = "UP"
				} else {
					upOrDown = "DOWN"
				}			  
				
				date := pastmatches[i].date.Month().String() + "-" + strconv.Itoa(pastmatches[i].date.Day()) + "-" + strconv.Itoa(pastmatches[i].date.Year())
				// Building the final output
				temp = temp + "#" + strconv.Itoa(i + 1) + " Against " + otherPlayer + " you played " + upOrDown + " and " + whoWon + " on " + date + "\n"
			}

			// Sending reponse if there is/isnt matches associated with the player
			if (len(temp) == 0) {
				response.Reply("Sorry, you don't have any stored matches")
			} else {
				response.Reply(temp)
			}
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

// updates the ladder after a match has been completed
func updateLadder(winner string, winnersOldPosition int, loser string, losersOldPosition int) {
	_, err := db.Exec("UPDATE rankings SET position='" + strconv.Itoa(winnersOldPosition - 1) + "'" + "WHERE firstname='" + winner + "'")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("UPDATE rankings SET position='" + strconv.Itoa(losersOldPosition + 1) + "'" + "WHERE firstname='" + loser + "'")
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated ladder!")
}

// records the match in the pastmatches table
func recordMatch(winner int, player1 string, player2 string, player1PrevPos int, player2PrevPos int) {
	player1Temp := "'" + player1 + "', "
	player2Temp := "'" + player2 + "', "
	winnerTemp := "'" + strconv.Itoa(winner) + "', "
	player1PrevPosTemp := "'" + strconv.Itoa(player1PrevPos) + "', "
	player2PrevPosTemp := "'" + strconv.Itoa(player2PrevPos) + "', "
	currentTime := "'" + time.Now().Format(time.RFC3339) + "'"
	_, err := db.Exec("INSERT INTO pastmatches VALUES (" + player1Temp + player2Temp + winnerTemp + player1PrevPosTemp + player2PrevPosTemp + currentTime + ")")
	if err != nil {
		panic(err)
	}
}

