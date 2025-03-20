package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"bowling-score-tracker/http_handlers"
)

func main() {
	r := gin.Default()
	http_handlers.RegisterEndpoints(r)

	if err := r.Run(":80"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
