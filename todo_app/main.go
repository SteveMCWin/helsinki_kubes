package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

var volume_path = "/usr/src/app/files/"
var image_name = "myimg.jpg"

func main() {

    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
	fmt.Println("Current directory: ", exPath)

	env_port := os.Getenv("PORT")
	if env_port == "" {
		env_port = "8080"
	}

	go func() {
		for {
			getAndStoreImage()
			time.Sleep(10 * time.Minute)
		}
	}()

	router := gin.Default()
	router.LoadHTMLGlob("front/*")
	router.Static("/files", "./files")

	router.GET("/", HandleGetHome())

	router.Run(":"+env_port)
}

func getAndStoreImage() {
	url := "https://picsum.photos/600/400"

    response, e := http.Get(url)
    if e != nil {
		panic(e)
    }
    defer response.Body.Close()

    //open a file for writing
    file, err := os.Create(filepath.Join(volume_path, image_name))
    if err != nil {
		panic(err)
    }
    defer file.Close()

    // Use io.Copy to just dump the response body to the file. This supports huge files
    _, err = io.Copy(file, response.Body)
    if err != nil {
		panic(err)
    }

	fmt.Println("Stored image!")
}

func HandleGetHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{ "photo": "/files/" + image_name })
	}
}
