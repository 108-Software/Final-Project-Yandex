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


// GetTaskHandler обработчик для GET /api/task
func GetTaskHandler(c *gin.Context) {
    // Получаем ID из query параметра
    idStr := c.Query("id")
    if idStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Не указан идентификатор",
        })
        return
    }
    
    // Конвертируем ID в int
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат идентификатора",
        })
        return
    }
    
    // Получаем задачу из БД
    task, err := db.GetTaskByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // Преобразуем в ответ со строковым ID
    taskResponse := TaskResponse{
        ID:      strconv.Itoa(task.ID),
        Date:    task.Date,
        Title:   task.Title,
        Comment: task.Comment,
        Repeat:  task.Repeat,
    }
    
    c.JSON(http.StatusOK, taskResponse)
}


// UpdateTaskHandler обработчик для PUT /api/task
func UpdateTaskHandler(c *gin.Context) {
    var taskReq TaskResponse
    
    // Парсим JSON из тела запроса
    if err := c.BindJSON(&taskReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат JSON",
        })
        return
    }
    
    // Проверяем обязательные поля
    if taskReq.ID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Не указан идентификатор",
        })
        return
    }
    
    if taskReq.Date == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Не указана дата",
        })
        return
    }
    
    if taskReq.Title == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Не указан заголовок",
        })
        return
    }
    
    // Проверяем валидность даты
    if !db.IsValidDate(taskReq.Date) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат даты",
        })
        return
    }
    
    // Проверяем длину заголовка (обычно есть ограничения)
    if len(taskReq.Title) > 255 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Слишком длинный заголовок",
        })
        return
    }
    
    // Конвертируем ID в int
    id, err := strconv.Atoi(taskReq.ID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат идентификатора",
        })
        return
    }
    
    // Создаем объект задачи для БД
    task := &db.Task{
        ID:      id,
        Date:    taskReq.Date,
        Title:   taskReq.Title,
        Comment: taskReq.Comment,
        Repeat:  taskReq.Repeat,
    }
    
    // Обновляем задачу в БД
    err = db.UpdateTask(task)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // Возвращаем пустой JSON при успехе
    c.JSON(http.StatusOK, gin.H{})
}
