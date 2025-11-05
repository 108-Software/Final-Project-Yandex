package db

import (
    "fmt"
)

func Tasks(limit int) ([]*Task, error) {
    if DB == nil {
        return nil, fmt.Errorf("база данных не инициализирована")
    }
    
    if limit <= 0 {
        limit = 50
    }
    
    query := `SELECT id, date, title, comment, repeat FROM scheduler 
              ORDER BY date ASC, id ASC 
              LIMIT ?`
    
    rows, err := DB.Query(query, limit)
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
