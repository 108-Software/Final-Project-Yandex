package api

import (
	"net/http"
	"time"
	"errors"
	"strconv"
	"strings"


	"github.com/gin-gonic/gin"
)

const dateFormat = "20060102"


func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	startDate, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", errors.New("некорректный формат исходной даты")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("неверный формат правила повторения")
	}

	command := parts[0]

	switch command {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат для правила 'd'")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("неверное число дней для правила 'd'")
		}

		if days <= 0 || days > 400 {
			return "", errors.New("интервал дней должен быть от 1 до 400")
		}

		return nextDateD(now, startDate, days), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат для правила 'y'")
		}

		return nextDateY(now, startDate), nil

	case "w":
		if len(parts) != 2 {
			return "", errors.New("неверный формат для правила 'w'")
		}

		dayStrs := strings.Split(parts[1], ",")
		days := make([]int, len(dayStrs))
		for i, dayStr := range dayStrs {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("недопустимое значение дня недели")
			}
			days[i] = day
		}

		return nextDateW(now, startDate, days), nil

	case "m":
		if len(parts) < 2 {
			return "", errors.New("неверный формат для правило 'm'")
		}

		dayStrs := strings.Split(parts[1], ",")
		days := make([]int, len(dayStrs))
		for i, dayStr := range dayStrs {
			day, err := strconv.Atoi(dayStr)
			if err != nil {
				return "", errors.New("неверный формат дня месяца")
			}
			if (day < -2 || day > 31) || day == 0 {
				return "", errors.New("недопустимый день месяца")
			}
			days[i] = day
		}

		var months []int
		if len(parts) >= 3 {
			monthStrs := strings.Split(parts[2], ",")
			months = make([]int, len(monthStrs))
			for i, monthStr := range monthStrs {
				month, err := strconv.Atoi(monthStr)
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("недопустимый месяц")
				}
				months[i] = month
			}
		}

		return nextDateMComplex(now, startDate, days, months), nil

	default:
		return "", errors.New("неизвестная команда правила повторения")
	}
}

func nextDateD(now time.Time, startDate time.Time, days int) string {
	currentDate := startDate

	for {
		currentDate = currentDate.AddDate(0, 0, days)
		if currentDate.After(now) {
			break
		}
	}

	return currentDate.Format(dateFormat)
}

func nextDateY(now time.Time, startDate time.Time) string {
	currentDate := startDate

	for {
		currentDate = currentDate.AddDate(1, 0, 0)
		if currentDate.After(now) {
			break
		}
	}

	return currentDate.Format(dateFormat)
}

func nextDateW(now time.Time, startDate time.Time, weekdays []int) string {
	currentDate := startDate

	for {
		currentDate = currentDate.AddDate(0, 0, 1)
		if !currentDate.After(now) {
			continue
		}

		currentWeekday := int(currentDate.Weekday())
		if currentWeekday == 0 {
			currentWeekday = 7
		}

		for _, targetWeekday := range weekdays {
			if currentWeekday == targetWeekday {
				return currentDate.Format(dateFormat)
			}
		}
	}
}

func nextDateMComplex(now time.Time, startDate time.Time, days, months []int) string {
	if len(months) == 0 {
		months = make([]int, 12)
		for i := 0; i < 12; i++ {
			months[i] = i + 1
		}
	}

	currentDate := startDate

	for {
		currentDate = currentDate.AddDate(0, 0, 1)
		if !currentDate.After(now) {
			continue
		}

		currentMonth := int(currentDate.Month())
		currentDay := currentDate.Day()

		monthMatch := false
		for _, month := range months {
			if currentMonth == month {
				monthMatch = true
				break
			}
		}

		if !monthMatch {
			continue
		}

		for _, day := range days {
			if day < 0 {
				lastDay := lastDayOfMonth(currentDate.Year(), time.Month(currentMonth), day)
				if currentDate.Equal(lastDay) {
					return currentDate.Format(dateFormat)
				}
			} else if currentDay == day {
				return currentDate.Format(dateFormat)
			}
		}
	}
}

func lastDayOfMonth(year int, month time.Month, dayFromEnd int) time.Time {
	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	return firstDayNextMonth.AddDate(0, 0, dayFromEnd)
}

func NextDateHandler(c *gin.Context) {
	nowStr := c.Query("now")
	date := c.Query("date")
	repeat := c.Query("repeat")

	if nowStr == "" || date == "" || repeat == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Отсутствуют обязательные параметры: now, date, repeat"})
		return
	}

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат параметра now"})
		return
	}

	nextDate, err := NextDate(now, date, repeat)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if nextDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Для указанного правила повторения нет подходящей даты"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"nextDate": nextDate})
}
