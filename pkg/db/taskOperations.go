package db

import (
    "database/sql"
    "fmt"
)


func GetTaskByID(id int) (*Task, error) {
    if DB == nil {
        return nil, fmt.Errorf("база данных не инициализирована")
    }
    
    query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
    
    var task Task
    err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("задача не найдена")
        }
        return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
    }
    
    return &task, nil
}


func UpdateTask(task *Task) error {
    if DB == nil {
        return fmt.Errorf("база данных не инициализирована")
    }
    
    _, err := GetTaskByID(task.ID)
    if err != nil {
        return err
    }
    
    if !IsValidDate(task.Date) {
        return fmt.Errorf("неверный формат даты")
    }
    
    if task.Title == "" {
        return fmt.Errorf("заголовок не может быть пустым")
    }
    
    if task.Repeat != "" && !IsValidRepeat(task.Repeat) {
        return fmt.Errorf("неверный формат повторения")
    }
    
    query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
    
    result, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
    if err != nil {
        return fmt.Errorf("ошибка выполнения запроса: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("ошибка проверки обновления: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("задача не найдена")
    }
    
    return nil
}
