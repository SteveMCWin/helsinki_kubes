package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type TodoItem struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var volume_path = "/usr/src/app/files/"
var image_name = "myimg.jpg"

func main() {

	all_args := os.Args
	if len(all_args) > 1 {
		volume_path = "./files/"
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println("Current directory: ", exPath)

	env_port := os.Getenv("PORT")
	if env_port == "" {
		env_port = ":8080"
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

	router.Run(env_port)
}

func getAndStoreImage() {

	//"https://picsum.photos" 
	env_pic_source := os.Getenv("PICTURE_SOURCE")
	env_width  := os.Getenv("PICTURE_WIDTH")
	env_height := os.Getenv("PICTURE_HEIGHT")

	url := env_pic_source + "/" + env_width + "/" + env_height

	// log.Println("URL: ", url)
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
		backend_url := os.Getenv("BACKEND_URL")
		res, err := http.Get(backend_url)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		var items []TodoItem

		err = json.NewDecoder(res.Body).Decode(&items)
		if err != nil {
			panic(err)
		}

		log.Println("Len of item list: ", len(items))

		c.HTML(http.StatusOK, "home.html", gin.H{"photo": "/files/" + image_name, "todoItems": items})
	}
}
