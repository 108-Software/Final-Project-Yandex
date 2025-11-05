package db

import (
    "fmt"
    "time"
)

func SearchTasks(limit int, search string) ([]*Task, error) {
    if DB == nil {
        return nil, fmt.Errorf("база данных не инициализирована")
    }
    
    if limit <= 0 {
        limit = 50
    }
    
    var query string
    var args []interface{}
    
    if search == "" {
        query = `SELECT id, date, title, comment, repeat FROM scheduler 
                 ORDER BY date ASC, id ASC 
                 LIMIT ?`
        args = []interface{}{limit}
    } else {
        if isDateSearch(search) {
            date, err := time.Parse("02.01.2006", search)
            if err == nil {
                searchDate := date.Format("20060102")
                query = `SELECT id, date, title, comment, repeat FROM scheduler 
                         WHERE date = ? 
                         ORDER BY date ASC, id ASC 
                         LIMIT ?`
                args = []interface{}{searchDate, limit}
            } else {
                query = `SELECT id, date, title, comment, repeat FROM scheduler 
                         WHERE title LIKE ? OR comment LIKE ? 
                         ORDER BY date ASC, id ASC 
                         LIMIT ?`
                searchPattern := "%" + search + "%"
                args = []interface{}{searchPattern, searchPattern, limit}
            }
        } else {
            query = `SELECT id, date, title, comment, repeat FROM scheduler 
                     WHERE title LIKE ? OR comment LIKE ? 
                     ORDER BY date ASC, id ASC 
                     LIMIT ?`
            searchPattern := "%" + search + "%"
            args = []interface{}{searchPattern, searchPattern, limit}
        }
    }
    
    rows, err := DB.Query(query, args...)
    if err != nil {
        return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
    }
    defer rows.Close()
    
    var tasks []*Task
    
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
        if err != nil {
            return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
        }
        tasks = append(tasks, &task)
    }
    
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
    }
    
    if tasks == nil {
        tasks = make([]*Task, 0)
    }
    
    return tasks, nil
}

func isDateSearch(s string) bool {
    if len(s) != 10 {
        return false
    }
    if s[2] != '.' || s[5] != '.' {
        return false
    }
    
    _, err := time.Parse("02.01.2006", s)
    return err == nil
}
