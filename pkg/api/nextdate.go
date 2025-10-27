package db

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	now := time.Date(2024, 1, 26, 0, 0, 0, 0, time.UTC)

	fmt.Println(NextDate(now, "20240229", "y"))
	fmt.Println(NextDate(now, "20240113", "d 7"))
	fmt.Println(NextDate(now, "20240116", "m 16,5"))
	fmt.Println(NextDate(now, "20240201", "m -1,18"))
	fmt.Println(NextDate(now, "20240101", "m -1"))
	fmt.Println(NextDate(now, "20240101", "m 3 1,3,6"))

	fmt.Println(NextDate(now, "20240101", "w 7"))
	fmt.Println(NextDate(now, "20240101", "w 1,4,5"))
	fmt.Println(NextDate(now, "20240101", "w 2,3"))
	fmt.Println(NextDate(now, "20240101", "w 6"))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	_, err := time.Parse("20060102", dstart)
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

		return nextDate_D(dstart, repeat), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("неверный формат для правила 'y'")
		}

		return nextDate_Y(dstart, repeat), nil

	case "w":
		if len(parts) != 2 {
			return "", errors.New("неверный формат для правила 'w'")
		}

		dayStrs := strings.Split(parts[1], ",")
		for _, dayStr := range dayStrs {
			day, err := strconv.Atoi(dayStr)
			if err != nil || day < 1 || day > 7 {
				return "", errors.New("недопустимое значение дня недели")
			}
		}

		return nextDate_W(now, dstart, repeat), nil

	case "m":
		if len(parts) < 2 {
			return "", errors.New("неверный формат для правила 'm'")
		}

		dayStrs := strings.Split(parts[1], ",")
		for _, dayStr := range dayStrs {
			day, err := strconv.Atoi(dayStr)
			if err != nil {
				return "", errors.New("неверный формат дня месяца")
			}
			if (day < -2 || day > 31) || day == 0 {
				return "", errors.New("недопустимый день месяца")
			}
		}

		if len(parts) >= 3 {
			monthStrs := strings.Split(parts[2], ",")
			for _, monthStr := range monthStrs {
				month, err := strconv.Atoi(monthStr)
				if err != nil || month < 1 || month > 12 {
					return "", errors.New("недопустимый месяц")
				}
			}
		}

		return nextDate_M_complex(now, dstart, repeat), nil

	default:
		return "", errors.New("неизвестная команда правила повторения")
	}
}

func nextDate_D(dstart, repeat string) string {
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		return "Ошибка даты"
	}

	numbers := repeat[2:]
	daysToAdd, _ := strconv.Atoi(numbers)

	resultDate := startDate.AddDate(0, 0, daysToAdd)
	return resultDate.Format("20060102")
}

func nextDate_Y(dstart, repeat string) string {
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		return "Ошибка даты"
	}

	resultDate := startDate.AddDate(1, 0, 0)
	return resultDate.Format("20060102")
}

func nextDate_M_complex(now time.Time, dstart, repeat string) string {
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		return "Ошибка даты"
	}

	parts := strings.Fields(repeat)
	if len(parts) < 2 {
		return "Ошибка: неверный формат правила"
	}

	daysStr := parts[1]

	var months []int
	if len(parts) >= 3 {
		monthStrs := strings.Split(parts[2], ",")
		for _, m := range monthStrs {
			month, _ := strconv.Atoi(m)
			if month >= 1 && month <= 12 {
				months = append(months, month)
			}
		}
	} else {
		for i := 1; i <= 12; i++ {
			months = append(months, i)
		}
	}

	dayStrs := strings.Split(daysStr, ",")
	var days []int
	for _, d := range dayStrs {
		day, _ := strconv.Atoi(d)
		days = append(days, day)
	}

	var possibleDates []time.Time
	startYear := startDate.Year()

	for year := startYear; year <= startYear+2; year++ {
		for _, month := range months {
			for _, day := range days {
				var date time.Time
				if day < 0 {
					date = lastDayOfMonth(year, time.Month(month), day)
					if date.Month() != time.Month(month) {
						continue
					}
				} else {
					date = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
					if date.Month() != time.Month(month) || date.Day() != day {
						continue
					}
				}

				if date.After(startDate) && date.After(now) {
					possibleDates = append(possibleDates, date)
				}
			}
		}
	}

	var closestDate time.Time
	found := false

	for _, date := range possibleDates {
		if !found || date.Before(closestDate) {
			closestDate = date
			found = true
		}
	}

	if !found {
		return "Нет подходящей даты"
	}

	return closestDate.Format("20060102")
}

func lastDayOfMonth(year int, month time.Month, dayFromEnd int) time.Time {
	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	return firstDayNextMonth.AddDate(0, 0, dayFromEnd)
}

func nextDate_W(now time.Time, dstart, repeat string) string {
	startDate, _ := time.Parse("20060102", dstart)

	parts := strings.Fields(repeat)
	daysStr := parts[1]

	dayStrs := strings.Split(daysStr, ",")
	var weekdays []int
	for _, dayStr := range dayStrs {
		day, _ := strconv.Atoi(dayStr)
		weekdays = append(weekdays, day)
	}

	currentDate := now.AddDate(0, 0, 1)

	for i := 0; i < 60; i++ {
		currentWeekday := int(currentDate.Weekday())
		if currentWeekday == 0 {
			currentWeekday = 7
		}

		for _, targetWeekday := range weekdays {
			if currentWeekday == targetWeekday {
				if currentDate.After(startDate) {
					return currentDate.Format("20060102")
				}
			}
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return "Нет подходящей даты"
}
