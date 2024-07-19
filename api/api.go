package api

import (
	"Task/api/handler"
	"Task/internal/mongodb"
	"Task/internal/rabbitmq"
	"log"


	"github.com/gin-gonic/gin"
)

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

	go rabbitmq.StartConsumer()

	router.Run(":7777")
}