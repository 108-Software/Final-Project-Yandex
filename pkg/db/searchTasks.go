package db

import (
    "fmt"
    "time"
)

// SearchTasks возвращает список задач с учетом поиска
// limit - максимальное количество возвращаемых записей
// search - строка поиска (может быть текстом или датой в формате 02.01.2006)
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
        // Если search пустой, работаем как обычная функция Tasks
        query = `SELECT id, date, title, comment, repeat FROM scheduler 
                 ORDER BY date ASC, id ASC 
                 LIMIT ?`
        args = []interface{}{limit}
    } else {
        // Проверяем, является ли search датой в формате 02.01.2006
        if isDateSearch(search) {
            // Преобразуем дату из 02.01.2006 в 20060102
            date, err := time.Parse("02.01.2006", search)
            if err == nil {
                searchDate := date.Format("20060102")
                query = `SELECT id, date, title, comment, repeat FROM scheduler 
                         WHERE date = ? 
                         ORDER BY date ASC, id ASC 
                         LIMIT ?`
                args = []interface{}{searchDate, limit}
            } else {
                // Если не удалось распарсить дату, ищем как текст
                query = `SELECT id, date, title, comment, repeat FROM scheduler 
                         WHERE title LIKE ? OR comment LIKE ? 
                         ORDER BY date ASC, id ASC 
                         LIMIT ?`
                searchPattern := "%" + search + "%"
                args = []interface{}{searchPattern, searchPattern, limit}
            }
        } else {
            // Поиск по тексту в заголовке или комментарии
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

// isDateSearch проверяет, является ли строка датой в формате 02.01.2006
func isDateSearch(s string) bool {
    // Проверяем базовую структуру: DD.MM.YYYY
    if len(s) != 10 {
        return false
    }
    if s[2] != '.' || s[5] != '.' {
        return false
    }
    
    // Пытаемся распарсить
    _, err := time.Parse("02.01.2006", s)
    return err == nil
}