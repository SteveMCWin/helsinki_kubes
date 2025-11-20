package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Read Postgres connection info from env vars
	pgHost := os.Getenv("DB_HOST")
	pgPort := os.Getenv("DB_PORT")
	pgUser := os.Getenv("DB_USER")
	pgPassword := os.Getenv("DB_PASSWORD")
	pgDB := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPassword, pgDB,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}
	defer db.Close()

	// Ensure the table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS counter (
		id SERIAL PRIMARY KEY,
		value INT NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// Insert initial row if table is empty
	var count int
	err = db.QueryRow("SELECT value FROM counter WHERE id = 1").Scan(&count)
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO counter (id, value) VALUES (1, 0)")
		if err != nil {
			log.Fatal("Failed to insert initial counter:", err)
		}
		count = 0
	} else if err != nil {
		log.Fatal("Failed to query counter:", err)
	}

	router := gin.Default()
	router.GET("/pings", func(c *gin.Context) {
		var val int
		err := db.QueryRow("SELECT value FROM counter WHERE id = 1").Scan(&val)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading counter")
			return
		}
		c.String(http.StatusOK, strconv.Itoa(val))
	})

	router.GET("/pingpong", func(c *gin.Context) {
		var val int
		err := db.QueryRow("SELECT value FROM counter WHERE id = 1").Scan(&val)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading counter")
			return
		}

		c.String(http.StatusOK, "pong "+strconv.Itoa(val))

		// Increment counter
		_, err = db.Exec("UPDATE counter SET value = value + 1 WHERE id = 1")
		if err != nil {
			log.Println("Failed to increment counter:", err)
		}
	})

	router.Run(":3000")
}
