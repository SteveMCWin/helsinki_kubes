package data

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type TodoItem struct {
	Id        int
	Name      string
	Completed bool
}

type Db struct {
	IsOpen bool
	Data   *sql.DB
}

func (db *Db) InitDb() {
	if db.IsOpen {
		log.Println("WARNING: Tried opening db that is already open")
		return
	}

	pgHost := os.Getenv("DB_HOST")
	pgPort := os.Getenv("DB_PORT")
	pgUser := os.Getenv("DB_USER")
	pgPassword := os.Getenv("DB_PASSWORD")
	pgDB := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost, pgPort, pgUser, pgPassword, pgDB,
	)

	var err error
	db.Data, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to Postgres:", err)
	}

	_, err = db.Data.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		completed BOOLEAN NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	db.IsOpen = true
}

func (db *Db) Close() {
	err := db.Data.Close()
	if err != nil {
		log.Println("Error closing database: ", err)
		return
	}

	db.IsOpen = false
}

func (db *Db) GetTodos() []TodoItem {
	rows, err := db.Data.Query("SELECT id, name, completed FROM todos ORDER BY id")
	if err != nil {
		log.Println("Error getting all tasks from todos: ", err)
		return nil
	}
	defer rows.Close()

	res := []TodoItem{}

	for rows.Next() {
		ti := TodoItem{}

		err = rows.Scan(&ti.Id, &ti.Name, &ti.Completed)
		if err != nil {
			log.Println("Error reading a rows: ", err)
			return nil
		}
		res = append(res, ti)
	}

	return res
}

func (db *Db) InsertTodo(todo *TodoItem) {
	statement := "INSERT into todos (name, completed) values ($1, $2) returning id"
	err := db.Data.QueryRow(statement, todo.Name, todo.Completed).Scan(&todo.Id)
	if err != nil {
		log.Println("Error inserting new todo into database: ", err)
	}
}

func (db *Db) UpdateTodo(todo *TodoItem) {
	statement := "UPDATE todos SET name = $2, completed = $3 where id = $1"
	_, err := db.Data.Exec(statement, todo.Id, todo.Name, todo.Completed)
	if err != nil {
		log.Println("Error inserting new todo into database: ", err)
	}
}
