package db

import (
    "time"
    "strings"
)

type Task struct {
    ID      int    `json:"id"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
    result, err := DB.Exec(
        `INSERT INTO scheduler (date, title, comment, repeat) 
         VALUES (?, ?, ?, ?)`,
        task.Date, task.Title, task.Comment, task.Repeat)
    if err != nil {
        return 0, err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }
    
    return id, nil
}

// IsValidDate проверяет валидность даты в формате 20060102
func IsValidDate(dateStr string) bool {
    if len(dateStr) != 8 {
        return false
    }
    
    // Проверяем, что все символы - цифры
    for _, char := range dateStr {
        if char < '0' || char > '9' {
            return false
        }
    }
    
    // Пытаемся распарсить
    _, err := time.Parse("20060102", dateStr)
    return err == nil
}

func IsValidRepeat(repeat string) bool {
    if repeat == "" {
        return true // пустое значение - допустимо
    }
    
    parts := strings.Fields(repeat)
    if len(parts) < 2 {
        return false
    }
    
    switch parts[0] {
    case "d": // ежедневно
        if len(parts) != 2 {
            return false
        }
        // Проверяем, что второй параметр - число
        for _, char := range parts[1] {
            if char < '0' || char > '9' {
                return false
            }
        }
        return true
        
    case "y", "w", "m": // ежегодно, еженедельно, ежемесячно
        return len(parts) == 1
        
    default:
        return false
    }
}
