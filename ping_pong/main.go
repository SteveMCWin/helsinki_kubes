package main

import (
	"net/http"
	// "os"
	// "path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// var volume_path = "/usr/src/app/files/"
// var pong_file_name = "pongs.txt"

func main() {
	counter := 0

	router := gin.Default()
	router.GET("/pings", HandleGetPings(&counter))
	router.GET("/pingpong", HandleGetHome(&counter))

	router.Run(":3000")
}

func HandleGetPings(counter *int) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, strconv.Itoa(*counter))
	}
}

func HandleGetHome(counter *int) func(c *gin.Context) {
	return func(c *gin.Context) {
		// to_write := strconv.Itoa(*counter)
		// err := os.WriteFile(filepath.Join(volume_path, pong_file_name), []byte(to_write), 0666)
		// if err != nil {
		// 	panic(err)
		// }
		c.String(http.StatusOK, "pong " + strconv.Itoa(*counter))
		(*counter)++
	}
}
