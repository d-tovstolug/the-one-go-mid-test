package controller

import (
	"errors"
	"github.com/d.tovstoluh/the-one-go-mid-test/rest-api/model"
	"github.com/d.tovstoluh/the-one-go-mid-test/rest-api/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	errInvalidRequest = errors.New("invalid request")
)

type TaskController struct {
	taskStorage storage.TaskStorage
}

func NewTaskController(taskStorage storage.TaskStorage) *TaskController {
	return &TaskController{
		taskStorage: taskStorage,
	}
}

func (c *TaskController) Register(r *gin.Engine) {
	taskRoutes := r.Group("v1/tasks")
	{
		taskRoutes.POST("/", c.SaveTask)
		taskRoutes.GET("/", c.GetAllTasks)
		taskRoutes.GET("/:id", c.GetTask)
		taskRoutes.DELETE("/:id", c.DeleteTask)
	}
}

// SaveTask
// @Summary Create or update task
// @Description Creates new or updates existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body model.Task true "Task data"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/tasks/ [post]
func (c *TaskController) SaveTask(ctx *gin.Context) {
	var task model.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errInvalidRequest.Error()})
		return
	}

	if task.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errInvalidRequest.Error()})
		return
	}

	newTask, err := c.taskStorage.Save(ctx, &task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, newTask)
}

// GetTask
// @Summary Retrieve a task by ID
// @Description Retrieves the details of a specific task by its ID
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string "ID is required"
// @Failure 404 {string} string "Task not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/tasks/{id} [get]
func (c *TaskController) GetTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	task, err := c.taskStorage.Get(ctx, id)
	switch {
	case errors.Is(err, storage.ErrNotFound):
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	case err != nil:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, task)
}

// GetAllTasks
// @Summary Retrieve all tasks
// @Description Retrieves a list of all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} model.Task
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/tasks/ [get]
func (c *TaskController) GetAllTasks(ctx *gin.Context) {
	tasks, err := c.taskStorage.GetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// DeleteTask
// @Summary Delete a task
// @Description Deletes a task by its ID
// @Tags tasks
// @Param id path string true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "ID is required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /v1/tasks/{id} [delete]
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	err := c.taskStorage.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
