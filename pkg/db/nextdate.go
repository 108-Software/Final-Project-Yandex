package db

import (
	"errors"
	"strconv"
	"strings"
	"time"
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
			if (day < -31 || day > 31) || day == 0 {
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
	current := startDate
	for {
		current = current.AddDate(0, 0, days)
		if current.After(now) {
			return current.Format(dateFormat)
		}
	}
}

func nextDateY(now time.Time, startDate time.Time) string {
	current := startDate
	for {
		current = current.AddDate(1, 0, 0)
		if current.After(now) {
			return current.Format(dateFormat)
		}
	}
}

func nextDateW(now time.Time, startDate time.Time, weekdays []int) string {
	current := now.AddDate(0, 0, 1)
	for i := 0; i < 365; i++ {
		currentWeekday := int(current.Weekday())
		if currentWeekday == 0 {
			currentWeekday = 7
		}
		for _, targetWeekday := range weekdays {
			if currentWeekday == targetWeekday {
				return current.Format(dateFormat)
			}
		}
		current = current.AddDate(0, 0, 1)
	}
	return ""
}

func nextDateMComplex(now time.Time, startDate time.Time, days, months []int) string {
	if len(months) == 0 {
		months = make([]int, 12)
		for i := 0; i < 12; i++ {
			months[i] = i + 1
		}
	}

	current := now
	for i := 0; i < 730; i++ {
		current = current.AddDate(0, 0, 1)
		currentMonth := int(current.Month())
		currentYear := current.Year()
		currentDay := current.Day()

		monthValid := false
		for _, m := range months {
			if currentMonth == m {
				monthValid = true
				break
			}
		}

		if monthValid {
			for _, day := range days {
				var checkDay int
				var dayExists bool

				if day > 0 {
					lastDay := lastDayOfMonth(currentYear, time.Month(currentMonth))
					if day > lastDay.Day() {
						continue
					}
					checkDay = day
					dayExists = true
				} else {
					lastDay := lastDayOfMonth(currentYear, time.Month(currentMonth))
					daysFromEnd := -day
					if daysFromEnd > lastDay.Day() {
						continue
					}
					checkDay = lastDay.Day() - daysFromEnd + 1
					dayExists = true
				}

				if dayExists && currentDay == checkDay && current.After(startDate) {
					return current.Format(dateFormat)
				}
			}
		}
	}
	return ""
}

func lastDayOfMonth(year int, month time.Month) time.Time {
	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	return firstDayNextMonth.AddDate(0, 0, -1)
}
