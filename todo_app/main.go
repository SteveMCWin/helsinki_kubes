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

type TodoItem struct {
	Name      string
	Completed bool
}

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

	todoItems := []TodoItem{
		TodoItem{Name: "Clean my room", Completed: true},
		TodoItem{Name: "Learn k8s"},
		TodoItem{Name: "Call mom"},
	}

	router.GET("/", HandleGetHome(todoItems))

	router.Run(":" + env_port)
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

func HandleGetHome(items []TodoItem) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{"photo": "/files/" + image_name, "todoItems": items})
	}
}
