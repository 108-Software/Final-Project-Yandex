package api

import (
    "net/http"
    "strconv"
    
    "Final-Project-Yandex/pkg/db"
    "github.com/gin-gonic/gin"
)

func TasksHandler(c *gin.Context) {
    search := c.Query("search")
    
    var tasks []*db.Task
    var err error
    
    if search == "" {
        tasks, err = db.Tasks(50)
    } else {
        tasks, err = db.SearchTasks(50, search)
    }
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Не удалось получить задачи: " + err.Error(),
        })
        return
    }
    
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


func GetTaskHandler(c *gin.Context) {
    idStr := c.Query("id")
    if idStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Не указан идентификатор",
        })
        return
    }
    
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат идентификатора",
        })
        return
    }
    
    task, err := db.GetTaskByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    taskResponse := TaskResponse{
        ID:      strconv.Itoa(task.ID),
        Date:    task.Date,
        Title:   task.Title,
        Comment: task.Comment,
        Repeat:  task.Repeat,
    }
    
    c.JSON(http.StatusOK, taskResponse)
}


func UpdateTaskHandler(c *gin.Context) {
    var taskReq TaskResponse
    
    if err := c.BindJSON(&taskReq); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат JSON",
        })
        return
    }
    
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
    
    if !db.IsValidDate(taskReq.Date) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат даты",
        })
        return
    }
    
    if len(taskReq.Title) > 255 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Слишком длинный заголовок",
        })
        return
    }
    
    id, err := strconv.Atoi(taskReq.ID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Неверный формат идентификатора",
        })
        return
    }
    
    task := &db.Task{
        ID:      id,
        Date:    taskReq.Date,
        Title:   taskReq.Title,
        Comment: taskReq.Comment,
        Repeat:  taskReq.Repeat,
    }
    
    err = db.UpdateTask(task)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{})
}
