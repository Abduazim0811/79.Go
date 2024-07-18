package api

import (
	"Task/api/handler"
	"Task/internal/mongodb"
	"Task/internal/rabbitmq"
	"log"

	ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
	_ "Task/cmd/docs"

	"github.com/gin-gonic/gin"
)
// New ...
// @title	Project: Swagger Intro
// @description This swagger UI was created in lesson
// @version 1.0
// @securityDefinitions.apikey Bearerauth
// @in header
// @name Authorization
func Routes(){
	router := gin.Default()
	db,err := mongodb.NewTask()
	if err != nil {
		log.Fatal(err)
	}
	taskhandler := handler.NewTaskHandler(db)

	router.POST("/tasks", taskhandler.CreateTask)
	router.GET("/tasks", taskhandler.GetTasks)
	router.PUT("/tasks/:id", taskhandler.UpdateTask)
	router.DELETE("/tasks/:id", taskhandler.DeleteTask)

	url := ginSwagger.URL("swagger/doc.json")
	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	go rabbitmq.StartConsumer()

	router.Run(":7777")
}