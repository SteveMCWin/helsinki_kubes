package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	counter := 0

	router := gin.Default()
	router.GET("/pingpong", HandleGetHome(&counter))

	router.Run(":3000")
}

func HandleGetHome(counter *int) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "pong " + strconv.Itoa(*counter))
		(*counter)++
	}
}
