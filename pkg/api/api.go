package api

import "github.com/gin-gonic/gin"

func Init(r *gin.Engine) {
	// Создаем группу для API маршрутов
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/nextdate", NextDateHandler)
		
	}
}