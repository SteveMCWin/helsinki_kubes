package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
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
	all_args := os.Args
	if len(all_args) < 2 {
		panic("Not enough arguments provided. Please pass 'read' or 'gen' as the arg to the executable")
	}

	args := all_args[1:]

	volume_path := "/usr/src/app/files"
	log_file_name := "log.txt"

	if args[0] == "gen" {
		genFunc(volume_path, log_file_name)
	} else {
		readFunc(volume_path, log_file_name)
	}
}

func genFunc(volume_path, log_file_name string) {
	random_str := randSeq(10)

	for {
		to_write := time.Now().Format(time.UnixDate) + "\t\t" + random_str
		err := os.WriteFile(filepath.Join(volume_path, log_file_name), []byte(to_write), 0666)
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func readFunc(volume_path, log_file_name string) {
	router := gin.Default()
	router.GET("/", HandleGetHome(volume_path, log_file_name))

	router.Run(":3000")
}

func HandleGetHome(volume_path, log_file_name string) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Println("Called handle get home")
		contents, err := os.ReadFile(filepath.Join(volume_path, log_file_name))
		if err != nil {
			panic(err)
		}

		c.String(http.StatusOK, string(contents))
	}
}
