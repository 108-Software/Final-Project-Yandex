package db

import (
    "time"
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

func IsValidDate(dateStr string) bool {
    _, err := time.Parse("20060102", dateStr)
    return err == nil
}
