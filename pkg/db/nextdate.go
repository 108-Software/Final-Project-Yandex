package db

import (
	"time"
	"errors"
	"strconv"
	"strings"
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
	// Начинаем с startDate и добавляем интервалы пока не получим дату после now
	current := startDate
	
	for {
		current = current.AddDate(0, 0, days)
		if current.After(now) {
			return current.Format(dateFormat)
		}
	}
}

func nextDateY(now time.Time, startDate time.Time) string {
	// Начинаем с startDate и добавляем годы пока не получим дату после now
	current := startDate
	
	for {
		nextYear := current.Year() + 1
		nextMonth := current.Month()
		nextDay := current.Day()
		
		// Проверяем валидность даты в следующем году
		nextDate := time.Date(nextYear, nextMonth, nextDay, 0, 0, 0, 0, time.UTC)
		
		// Если дата невалидна (например, 29 февраля в невисокосном году),
		// используем 1 марта следующего года
		if nextDate.Month() != nextMonth || nextDate.Day() != nextDay {
			nextDate = time.Date(nextYear, 3, 1, 0, 0, 0, 0, time.UTC)
		}
		
		current = nextDate
		if current.After(now) {
			return current.Format(dateFormat)
		}
	}
}

func nextDateW(now time.Time, startDate time.Time, weekdays []int) string {
	// Начинаем поиск со дня после now
	current := now.AddDate(0, 0, 1)
	
	// Ищем ближайший подходящий день недели
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
	
	return "" // не нашли подходящую дату
}

func nextDateMComplex(now time.Time, startDate time.Time, days, months []int) string {
	if len(months) == 0 {
		months = make([]int, 12)
		for i := 0; i < 12; i++ {
			months[i] = i + 1
		}
	}

	// Начинаем поиск с дня после now
	current := now.AddDate(0, 0, 1)
	
	// Ищем ближайшую подходящую дату в пределах 2 лет
	for i := 0; i < 730; i++ {
		currentMonth := int(current.Month())
		currentDay := current.Day()

		// Проверяем месяц
		monthMatch := false
		for _, month := range months {
			if currentMonth == month {
				monthMatch = true
				break
			}
		}

		if monthMatch {
			// Проверяем день
			for _, day := range days {
				if day < 0 {
					// Отрицательные дни - с конца месяца
					lastDay := lastDayOfMonth(current.Year(), time.Month(currentMonth))
					daysFromEnd := lastDay.Day() - currentDay + 1
					if -day == daysFromEnd {
						return current.Format(dateFormat)
					}
				} else if currentDay == day {
					// Положительные дни
					return current.Format(dateFormat)
				}
			}
		}
		current = current.AddDate(0, 0, 1)
	}
	
	return "" // не нашли подходящую дату
}

func lastDayOfMonth(year int, month time.Month) time.Time {
	firstDayNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	return firstDayNextMonth.AddDate(0, 0, -1)
}