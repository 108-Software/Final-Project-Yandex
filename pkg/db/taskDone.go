package db

import (
    "fmt"
    "time"
)

func MarkTaskDone(id int, now time.Time, nextDateFunc func(time.Time, string, string) (string, error)) error {
    if DB == nil {
        return fmt.Errorf("база данных не инициализирована")
    }
    
    task, err := GetTaskByID(id)
    if err != nil {
        return err
    }
    
    if task.Repeat == "" {
        return DeleteTask(id)
    }
    
    nextDate, err := nextDateFunc(now, task.Date, task.Repeat)
    if err != nil {
        return fmt.Errorf("не удалось рассчитать следующую дату: %v", err)
    }
    
    if nextDate == "" {
        return fmt.Errorf("не удалось рассчитать следующую дату")
    }
    
    return UpdateTaskDate(id, nextDate)
}

func DeleteTask(id int) error {
    query := `DELETE FROM scheduler WHERE id = ?`
    result, err := DB.Exec(query, id)
    if err != nil {
        return fmt.Errorf("ошибка удаления задачи: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("ошибка проверки удаления: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("задача не найдена")
    }
    
    return nil
}

func UpdateTaskDate(id int, newDate string) error {
    query := `UPDATE scheduler SET date = ? WHERE id = ?`
    result, err := DB.Exec(query, newDate, id)
    if err != nil {
        return fmt.Errorf("ошибка обновления даты: %v", err)
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