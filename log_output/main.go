package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}

func main() {
	random_str := randSeq(10)

	router := gin.Default()
	router.GET("/", HandleGetHome(random_str))

	router.Run(":3000")
}

func HandleGetHome(status string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, time.Now().Format(time.UnixDate) + "\t\t" + status)
	}
}
