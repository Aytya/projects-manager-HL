package handler

import (
	"database/sql"
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateProject
// @Summary      create project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        input body entity.Project true "Project Entity"
// @Success      201  {object}	entity.Project
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /projects [post]
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

// UpdateProject
// @Summary      update project by id
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        id path string true "Project ID"
// @Param        input body entity.Project true "Project Entity"
// @Success      200  {string}  ""message": "Task updated""
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /projects/{id} [put]
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

// DeleteProject
// @Summary      delete project by id
// @Tags         projects
// @Produce      json
// @Param        id path string true "Project ID"
// @Success      200  {object}	entity.Project
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /projects/{id} [delete]
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

// GetProjectById
// @Summary      get project by id
// @Tags         projects
// @Produce      json
// @Param        id path string true "Project ID"
// @Success      200  {object}	entity.Project
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /projects/{id} [get]
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

// GetAllProjects
// @Summary      get all projects
// @Tags         projects
// @Produce      json
// @Success      200  {array}	entity.Project
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /projects [get]
func (h *Handler) getListOfProjects(c *gin.Context) {
	lists, err := h.service.GetAllProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": lists})
}

// GetProjectTasks
// @Summary      get project's tasks
// @Tags         projects
// @Produce      json
// @Param        id path string true "Project ID"
// @Success      200  {object}	entity.Task
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /projects/{id}/tasks [get]
func (h *Handler) getProjectTasks(c *gin.Context) {
	h.getTaskByProjectId(c, "id")
}

// GetProjectByTitle
// @Summary      get project by title
// @Tags         projects
// @Produce      json
// @Param        title path string true "Project Title"
// @Success      200  {object}	entity.Project
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /projects/search/{title} [get]
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

// GetProjectByManagerId
// @Summary      get project by managerId
// @Tags         projects
// @Produce      json
// @Param        manager query string true "Project Manager"
// @Success      200  {object}	entity.Project
// @Failure      400  {object}  response.Object
// @Failure      404  {object}  response.Object
// @Failure      500  {object}  response.Object
// @Router       /projects/search [get]
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
