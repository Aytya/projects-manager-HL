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

func (h *Handler) getTasksByProjectId(c *gin.Context) {
	h.getTaskByProjectId(c, taskProject)
}

func (h *Handler) getAllTasks(c *gin.Context) {
	lists, err := h.service.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all tasks", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": lists})
}

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
