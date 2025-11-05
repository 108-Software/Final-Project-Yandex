package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"Final-Project-Yandex/pkg/db"
)

func AddTaskHandler(c *gin.Context) {
	var req struct {
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка десериализации JSON"})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан заголовок задачи"})
		return
	}

	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	if err := checkDate(task); err != nil {
		switch err.(type) {
		case *time.ParseError:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	id, err := db.AddTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении задачи в базу данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": strconv.FormatInt(id, 10)})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)
	
	if task.Date == "today" {
		task.Date = today
		return nil
	}
	
	if task.Date == "" {
		task.Date = today
		return nil
	}
	
	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}
	
	tNormalized := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	nowNormalized := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	if tNormalized.Before(nowNormalized) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			next, err := NextDate(now, today, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
	}
	
	return nil
}
