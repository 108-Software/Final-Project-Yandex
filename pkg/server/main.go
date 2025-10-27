package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	_ "modernc.org/sqlite"

	"Final-Project-Yandex/pkg/db"
)

func main() {

	dbPath := "./pkg/db/scheduler.db"
	if envPath := os.Getenv("TODO_DBFILE"); envPath != "" {
		dbPath = envPath
	}
	log.Printf("Используется файл базы данных: %s", dbPath)

	err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer db.Close()

	database := db.GetDB()
	if database == nil {
		log.Fatal("Соединение с базой данных не установлено")
	}

	log.Println("\nБаза данных доступна\n")

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

