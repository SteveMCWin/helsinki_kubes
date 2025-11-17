package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type TodoItem struct {
	Name      string
	Completed bool
}

var todoItems []TodoItem

func main() {

	env_port := os.Getenv("BACKENDPORT")
	if env_port == "" {
		env_port = "8080"
	}

	router := gin.Default()

	todoItems = []TodoItem{
		TodoItem{Name: "Clean my room", Completed: true},
		TodoItem{Name: "Learn k8s"},
		TodoItem{Name: "Call mom"},
	}

	router.GET("/todos", HandleGetTodos())
	router.POST("/todos", HandlePostTodos())

	router.Run(":" + env_port)
}

func HandleGetTodos() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, todoItems)
	}
}

func HandlePostTodos() func(c *gin.Context) {
	return func(c *gin.Context) {
		todo := c.PostForm("todo")

		if todo == "" {
			c.JSON(400, gin.H{"error": "todo cannot be empty"})
			return
		}

		log.Println("Before adding: ", todoItems)
		todoItems = append(todoItems, TodoItem{Name: todo})
		log.Println("After adding:  ", todoItems)

		c.JSON(200, gin.H{
			"status": "ok",
			"todo":   todo,
		})
	}
}
