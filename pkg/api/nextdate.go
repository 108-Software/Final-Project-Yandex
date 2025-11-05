package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"Final-Project-Yandex/pkg/db"
)

const dateFormat = "20060102"


func NextDateHandler(c *gin.Context) {
	nowStr := c.Query("now")
	date := c.Query("date")
	repeat := c.Query("repeat")

	if date == "" {
		c.String(400, "Параметр 'date' обязателен")
		return
	}
	if repeat == "" {
		c.String(400, "Параметр 'repeat' обязателен")
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			c.String(400, "Неверный формат параметра 'now'")
			return
		}
	}

	nextDate, err := db.NextDate(now, date, repeat)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	c.String(200, nextDate)
}
