package handler

import (
	"Task/internal/models"
	"Task/internal/mongodb"
	"Task/internal/rabbitmq"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	db *mongodb.TaskMongoDb
}

func NewTaskHandler(db *mongodb.TaskMongoDb) *TaskHandler {
	return &TaskHandler{db: db}
}


// @Router 			/tasks [post]
// @Summary			CREATE TASKS
// @Description 	This method creates tasks
// @Security		BearerAuth
// @Tags			TASK
// @Accept 			json
// @Produce			json
// @Param 			body body models.Tasks true "Tasks"
// @Success 201 {object} models.Tasks
// @Failure 400 {object} gin.H{"error": "Invalid input"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /tasks [post]
func (t *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Tasks
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.Id = primitive.NewObjectID()
	task.Status = "pending"
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	if err := rabbitmq.PublishTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}
// GetTasks godoc
// @Summary Get all tasks
// @Description Get all tasks
// @Tags tasks
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Tasks
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /tasks [get]
func (t *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := t.db.StoreGetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update a task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path string true "Tasks ID"
// @Param task body models.Tasks true "Tasks"
// @Success 200 {object} models.Task
// @Failure 400 {object} gin.H{"error": "Invalid input"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /tasks/{id} [put]
func (t *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Tasks
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.UpdatedAt = time.Now()

	if err = t.db.StoreUpdateTask(objectId, task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param id path string true "cdTask ID"
// @Success 204
// @Failure 400 {object} gin.H{"error": "Invalid task ID"}
// @Failure 500 {object} gin.H{"error": "Internal server error"}
// @Router /tasks/{id} [delete]
func (t *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := t.db.StoreDeleteTask(objectId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
