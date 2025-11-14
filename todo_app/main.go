package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	env_port := os.Getenv("PORT")
	if env_port == "" {
		env_port = "8080"
	}

	router := gin.Default()
	router.LoadHTMLGlob("front/*")

	router.GET("/", HandleGetHome())

	router.Run(":"+env_port)
}

func HandleGetHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	}
}
