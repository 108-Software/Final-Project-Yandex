package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	_ "modernc.org/sqlite"

	"Final-Project-Yandex/pkg/db"
	"Final-Project-Yandex/pkg/api"
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

	log.Println("База данных доступна")

	r := gin.Default()
	
	// Регистрируем API маршруты
	api.Init(r)
	
	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")
	r.Static("/images", "./web/images")
	
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/index.html", "./web/index.html")
	
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil && p > 0 && p < 65536 {
			port = envPort
		}
	}
	
	log.Printf("Сервер запущен на http://localhost:%s", port)	
	r.Run(":" + port)
}
