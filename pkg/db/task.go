package db

import (
	"time"
)

// Task структура задачи
type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в базу данных
func AddTask(task *Task) (int64, error) {
	// Выполняем SQL запрос
	result, err := db.Exec(`
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (?, ?, ?, ?)
	`, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	
	// Получаем ID последней вставленной записи
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

// Вспомогательная функция для проверки даты
func IsValidDate(dateStr string) bool {
	_, err := time.Parse("20060102", dateStr)
	return err == nil
}