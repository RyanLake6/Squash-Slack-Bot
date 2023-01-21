package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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
				response.Reply("Great job " + player1Name + "! \nThe ladder won't change")
			} else if (int(player1Position) == int(player2Position) + 1 && player1Position > player2Position) {
				response.Reply("Great job " + player1Name + "! \nThe ladder will be affected")
				updateLadder(player1Name, int(player1Position), player2Name, int(player2Position))
			} else {
				response.Reply("uh oh I can't support players beating eachother not next to eachother in the ladder :(")
			}
		},
	})

	// bot.Command("randomly select someone", &slacker.CommandDefinition{
	// 	Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
	// 		rows, _ := db.Query("SELECT firsName from rankings ORDER by RAND() LIMIT 1")
	// 		player := []Player{}
	// 		defer rows.Close()
	// 		for rows.Next() {
	// 			var p Player
	// 			rows.Scan(&p.position, &p.firstName)
	// 			player = append(player, p)
	// 		}
	// 		response.Reply(player[0].firstName + " you have been selected")
	// 	},
	// })

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

