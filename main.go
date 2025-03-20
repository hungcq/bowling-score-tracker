package main

import (
	"bowling-score-tracker/http_handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	http_handlers.RegisterEndpoints(r)

	if err := r.Run(":80"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
