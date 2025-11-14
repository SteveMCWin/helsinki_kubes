package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	env_port := os.Getenv("PORT")
	if env_port == "" {
		env_port = "8080"
	}

	router := gin.Default()

	router.Run(":"+env_port)
}
