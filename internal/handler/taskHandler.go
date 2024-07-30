package handler

import (
	"fmt"
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

const (
	taskId       = "id"
	taskTitle    = "title"
	taskStatus   = "status"
	taskPriority = "priority"
	taskProject  = "project"
)

// CreateTask
// @Summary      create task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        input body entity.Task true "Task Entity"
// @Success      201  {object}	entity.Task
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks [post]
func (h *Handler) createTask(c *gin.Context) {
	var task entity.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if task.Title == "" || task.Description == "" || task.Priority == "" || task.Status == "" || task.Assignee == "" || task.Project == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	id, createdAt, err := h.service.Task.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "created_at": createdAt})
}

// GetTaskById
// @Summary      get task by id
// @Tags         tasks
// @Produce      json
// @Param        id path string true "Task ID"
// @Success      200  {object}	entity.Task
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/{id} [get]
func (h *Handler) getTaskById(c *gin.Context) {
	id := c.Param(taskId)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	task, err := h.service.GetTaskById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})

}

// GetTaskByTitle
// @Summary      get task by title
// @Tags         tasks
// @Produce      json
// @Param        title query string true "Title"
// @Success      200  {object}	entity.Task
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/search [get]
func (h *Handler) getTaskByTitle(c *gin.Context) {
	title := c.Query(taskTitle)
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	task, err := h.service.GetTaskByTitle(title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// GetTasksByStatus
// @Summary      get tasks by status
// @Tags         tasks
// @Produce      json
// @Param        status query string true "Status"
// @Success      200  {object}	entity.Task
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/search/status [get]
func (h *Handler) getTaskByStatus(c *gin.Context) {
	status := c.Query(taskStatus)
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}

	task, err := h.service.GetTaskByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

// GetTasksByPriority
// @Summary      get tasks by priority
// @Tags         tasks
// @Produce      json
// @Param        priority query string true "Priority"
// @Success      200  {object}	entity.Task
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/search/priority [get]
func (h *Handler) getTaskByPriority(c *gin.Context) {
	priority := c.Query(taskPriority)
	if priority == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "priority is required"})
		return
	}

	task, err := h.service.GetTaskByPriority(priority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (h *Handler) getTasksByUserId(c *gin.Context, column string, value string) {
	assignee := c.Param(value)
	if assignee == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", value)})
		return
	}

	if _, err := uuid.Parse(assignee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format for assignee"})
		return
	}

	tasks, err := h.service.GetTasksByUserId(assignee, column)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})

}

// GetTasksByAssignee
// @Summary      get tasks by assignee id
// @Tags         tasks
// @Produce      json
// @Param        assignee path string true "Assignee ID"
// @Success      200  {object}	entity.Task
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/search/{assignee} [get]
func (h *Handler) getTasksByAssignee(c *gin.Context) {
	h.getTasksByUserId(c, "assignee", "assignee")
}

func (h *Handler) getTaskByProjectId(c *gin.Context, paramName string) {
	projectId := c.Param(paramName)
	if projectId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s is required", paramName)})
		return
	}

	tasks, err := h.service.GetTasksByProjectId(projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetTasksByProjectId
// @Summary      get tasks by project id
// @Tags         tasks
// @Produce      json
// @Param        project path string true "Project ID"
// @Success      200  {object}	entity.Task
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/search/project/{project} [get]
func (h *Handler) getTasksByProjectId(c *gin.Context) {
	h.getTaskByProjectId(c, taskProject)
}

// GetAllTasks
// @Summary      get all tasks
// @Tags         tasks
// @Produce      json
// @Success      200  {array}	entity.Task
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks [get]
func (h *Handler) getAllTasks(c *gin.Context) {
	lists, err := h.service.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all tasks", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": lists})
}

// UpdateTask
// @Summary      update task by id
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Task ID"
// @Param        input body entity.Task true "Project Entity"
// @Success      200  {string}  ""message": "Task updated""
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/{id} [put]
func (h *Handler) updateTaskById(c *gin.Context) {
	id := c.Param(taskId)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}

	var input entity.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := h.service.UpdateTask(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task updated"})
}

// DeleteTask
// @Summary      delete task by id
// @Tags         tasks
// @Produce      json
// @Param        id path string true "Task ID"
// @Success      200  "No content"
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /tasks/{id} [delete]
func (h *Handler) deleteTaskById(c *gin.Context) {
	id := c.Param(taskId)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}

	if err := h.service.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task", "message": err.Error()})
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
