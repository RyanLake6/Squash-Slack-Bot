package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func DatabaseConnect() *sql.DB {
	// Waiting for the db container to spin up??
	time.Sleep(time.Minute * 2)

	godotenv.Load()
	// Capture connection properties.
	cfg := mysql.Config{
		ParseTime: true,
		User:      "my_user", //os.Getenv("DBUSER"),
		Passwd:    "my_password", // os.Getenv("DBPASS"),
		Net:       "tcp",
		Addr:      "db:3306",
		DBName:    "ladder",
	}
	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to sql!")
	return db
}