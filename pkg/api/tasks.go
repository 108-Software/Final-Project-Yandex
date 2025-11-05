package api

import (
    "net/http"
    "strconv"
    
    "Final-Project-Yandex/pkg/db"
    "github.com/gin-gonic/gin"
)

func TasksHandler(c *gin.Context) {
    // Получаем параметр search из query string
    search := c.Query("search")
    
    var tasks []*db.Task
    var err error
    
    if search == "" {
        // Используем старую функцию если search пустой
        tasks, err = db.Tasks(50)
    } else {
        // Используем новую функцию поиска
        tasks, err = db.SearchTasks(50, search)
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Не удалось получить задачи: " + err.Error(),
        })
        return
    }
    
    // Преобразуем задачи из db.Task в TaskResponse со строковыми ID
    taskResponses := make([]*TaskResponse, 0, len(tasks))
    for _, task := range tasks {
        taskResponses = append(taskResponses, &TaskResponse{
            ID:      strconv.Itoa(task.ID),
            Date:    task.Date,
            Title:   task.Title,
            Comment: task.Comment,
            Repeat:  task.Repeat,
        })
    }
    
    c.JSON(http.StatusOK, gin.H{
        "tasks": taskResponses,
    })
}