package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	"Final-Project-Yandex/pkg/api"
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

	log.Println("База данных доступна")

	r := gin.Default()

	r.Static("/css", "./web/css")
	r.Static("/js", "./web/js")
	r.Static("/images", "./web/images")
	r.StaticFile("/favicon.ico", "./web/favicon.ico")

	api.Init(r)

	r.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	r.GET("/index.html", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	r.GET("/login.html", func(c *gin.Context) {
		c.File("./web/login.html")
	})

	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil && p > 0 && p < 65536 {
			port = envPort
		}
	}

	log.Printf("Сервер запущен на http://localhost:%s", port)
	r.Run(":" + port)
}
