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
		log.Fatal("Failed to create database handler: ", err)
	}

	maxRetries := 10
	retryInterval := 3 * time.Second
	for i := 0; i < maxRetries; i++ {
		pingErr := db.Ping()
		if err == nil {
			log.Printf("Connected to Database!")
			break
		}

		log.Printf("Attempt %d: Failed to connect to database: %s\n", i+1, err.Error())

		if i < maxRetries-1 {
			log.Printf("Retrying in %s...\n", retryInterval)
			time.Sleep(retryInterval)
		} else {
			log.Fatal(pingErr)
		}
	}
	

	fmt.Println("Connected to sql!")
	return db
}