package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"todo_backend/data"

	"github.com/gin-gonic/gin"
)

const MAX_TODO_LEN = 140

func main() {

	var db data.Db
	db.InitDb()

	env_port := os.Getenv("BACKENDPORT")
	if env_port == "" {
		env_port = ":8080"
	}

	router := gin.Default()

	router.GET("/todos", HandleGetTodos(&db))
	router.POST("/todos", HandlePostTodos(&db))
	router.PUT("/todos", HandlePutTodos(&db))

	router.Run(env_port)
}

func HandleGetTodos(db *data.Db) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Println("HandleGetTodos called.")
		todos, err := db.GetTodos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, todos)
	}
}

func HandlePostTodos(db *data.Db) func(c *gin.Context) {
	return func(c *gin.Context) {
		todo := c.PostForm("todo")
		log.Println("HandlePostTodos called. The todo passed in is: ", todo)

		if len(todo) > MAX_TODO_LEN {
			message := "todo name must not be longer than " + strconv.Itoa(MAX_TODO_LEN) + " characters"
			log.Println("Error: ", message)
			c.JSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}

		if todo == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "todo cannot be empty"})
			return
		}

		todoItem := data.TodoItem{Name: todo, Completed: false}
		err := db.InsertTodo(&todoItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"todo":   todoItem,
		})
	}
}

func HandlePutTodos(db *data.Db) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Println("HandlePutTodos called.")
		decoder := json.NewDecoder(c.Request.Body)
		var todo data.TodoItem
		err := decoder.Decode(&todo)
		if err != nil {
			log.Println("Error decoding todo: ", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		if len(todo.Name) > MAX_TODO_LEN {
			message := "todo name must not be longer than " + strconv.Itoa(MAX_TODO_LEN) + " characters"
			log.Println("Error: ", message)
			c.JSON(http.StatusBadRequest, gin.H{"error": message})
			return
		}

		err = db.UpdateTodo(&todo)
		if err != nil {
			log.Println("Error decoding todo: ", err)
			c.JSON(http.StatusInternalServerError, nil)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"todo": todo,
		})

	}
}
