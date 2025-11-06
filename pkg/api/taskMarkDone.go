package api

import (
    "net/http"
    "strconv"
    "time"
	"fmt"
    "Final-Project-Yandex/pkg/db"
    "github.com/gin-gonic/gin"
)

func TaskDoneHandler(c *gin.Context) {
    idStr := c.Query("id")
    if idStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан идентификатор"})
        return
    }
    
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат идентификатора"})
        return
    }
    
    
    err = markTaskDone(id, time.Now())
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{})
}

func DeleteTaskHandler(c *gin.Context) {
    
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
    
    err = db.DeleteTask(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{})
}


func markTaskDone(id int, now time.Time) error {
    task, err := db.GetTaskByID(id)
    if err != nil {
        return err
    }
    
    if task.Repeat == "" {
        return db.DeleteTask(id)
    }
    
    nextDate, err := db.NextDate(now, task.Date, task.Repeat)
    if err != nil {
        return err
    }
    
    if nextDate == "" {
        return fmt.Errorf("не удалось рассчитать следующую дату")
    }
    
    return db.UpdateTaskDate(id, nextDate)
}
