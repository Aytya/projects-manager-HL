package handler

import (
	"database/sql"
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) createProject(c *gin.Context) {
	var project entity.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	id, createdAt, err := h.service.Project.CreateProject(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        id,
		"createdAt": createdAt,
	})
}

func (h *Handler) updateProject(c *gin.Context) {
	id := c.Param("id")
	var updatedProject entity.Project
	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateProject(id, updatedProject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated"})
}

func (h *Handler) deleteProject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}

	if err := h.service.DeleteProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
}

func (h *Handler) getProjectById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
	}

	project, err := h.service.GetProjectById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

func (h *Handler) getListOfProjects(c *gin.Context) {
	lists, err := h.service.GetAllProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": lists})
}

func (h *Handler) getProjectTasks(c *gin.Context) {
	h.getTaskByProjectId(c, "id")
}

func (h *Handler) getProjectByTitle(c *gin.Context) {
	title := c.Param("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
	}

	project, err := h.service.GetProjectByTitle(title)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

func (h *Handler) getProjectByManagerId(c *gin.Context) {
	id := c.Query("manager")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "managerId is required"})

	}

	project, err := h.service.GetProjectByManagerId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}
