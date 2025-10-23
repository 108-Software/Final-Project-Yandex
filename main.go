package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

func main() {
	r := gin.Default()
	r.Static("/", "./web")
	
	port := "7540" // значение по умолчанию
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil && p > 0 && p < 65536 {
			port = envPort
		}
	}
	
	log.Printf("Сервер запущен на http://localhost:%s", port)
	r.Run(":" + port)
}