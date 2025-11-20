package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	// "path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var volume_path = "/usr/src/app/files/"
var config_file_name = "information.txt"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var genContents string

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

	if args[0] == "gen" {
		genFunc()
	} else {
		go func() {
			random_str := randSeq(10)

			for {
				// to_write := time.Now().Format(time.UnixDate) + "\t\t" + random_str
				genContents = time.Now().Format(time.UnixDate) + "\t\t" + random_str
				// err := os.WriteFile(filepath.Join(volume_path, log_file_name), []byte(to_write), 0666)
				// if err != nil {
				// 	panic(err)
				// }
				time.Sleep(5 * time.Second)
			}
		}()
		readFunc()
	}
}

func genFunc() {
	random_str := randSeq(10)

	for {
		// to_write := time.Now().Format(time.UnixDate) + "\t\t" + random_str
		genContents = time.Now().Format(time.UnixDate) + "\t\t" + random_str
		// err := os.WriteFile(filepath.Join(volume_path, log_file_name), []byte(to_write), 0666)
		// if err != nil {
		// 	panic(err)
		// }
		time.Sleep(5 * time.Second)
	}
}

func readFunc() {
	router := gin.Default()
	router.GET("/", HandleGetHome())

	router.Run(":3000")
}

func HandleGetHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Println("Called handle get home")
		config_contents, err := os.ReadFile(filepath.Join(volume_path, config_file_name))
		if err != nil {
			panic(err)
		}

		config_msg := os.Getenv("MESSAGE")

		res, err := http.Get("http://pingpong-svc:3456/pings")
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		resBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		res_string := 
			"this is read from a file: " + string(config_contents) +
			"\nenv variable: " + config_msg +
			"\n" + genContents + 
			"\nPing pongs: " + string(resBytes)

		c.String(http.StatusOK, res_string)
	}
}
