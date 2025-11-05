package api

import "github.com/gin-gonic/gin"

type TaskResponse struct {
    ID      string `json:"id"`
    Date    string `json:"date"`
    Title   string `json:"title"`
    Comment string `json:"comment"`
    Repeat  string `json:"repeat"`
}

type TasksResp struct {
    Tasks []*TaskResponse `json:"tasks"`
}

func Init(r *gin.Engine) {
	//Группа для API маршрутов
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/nextdate", NextDateHandler)
		apiGroup.GET("/tasks", TasksHandler)
		apiGroup.POST("/task", AddTaskHandler)
	}

}
