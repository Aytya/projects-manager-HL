package handler

import (
	"database/sql"
	"errors"
	"github.com/Aytya/projects-manager-HL/internal/entity"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	userId    = "id"
	userName  = "name"
	userEmail = "email"
)

var projects []entity.Project

func (h *Handler) createUser(c *gin.Context) {
	var newUser entity.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "message": err.Error()})
		return
	}

	id, createdAt, err := h.service.User.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        id,
		"createdAt": createdAt,
	})
}

func (h *Handler) getTasksByUser(c *gin.Context) {
	h.getTasksByUserId(c, "assignee", "id")
}

func (h *Handler) getUserById(c *gin.Context) {
	id := c.Param(userId)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := h.service.User.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getUser(c *gin.Context, queryParam, errorMessage string, getUserFunc func(string) (entity.User, error)) {
	value := c.Query(queryParam)
	if value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	user, err := getUserFunc(value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getUserByName(c *gin.Context) {
	h.getUser(c, userName, "Name query parameter is required", h.service.User.GetUserByName)
}

func (h *Handler) getUserByEmail(c *gin.Context) {
	h.getUser(c, userEmail, "Email query parameter is required", h.service.User.GetUserByEmail)
}

func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param(userId)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var input entity.User
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input", "message": err.Error()})
		return
	}

	if err := h.service.User.UpdateUser(id, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *Handler) deleteUser(c *gin.Context) {
	id := c.Param(userId)
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	if err := h.service.User.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *Handler) getAllUsers(c *gin.Context) {
	lists, err := h.service.User.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lists)
}
