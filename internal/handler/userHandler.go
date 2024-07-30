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

// Create User
// @Summary create user
// @Tags users
// @Accept json
// @Product json
// @Param		request	body		entity.User	true	"body param"
// @Success	200		{object}	entity.User
// @Failure	404	{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router       /users [post]
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

// @Summary get tasks by user id
// @Tags users
// @Product json
// @Param id path string true "user id"
// @Success	200		{object}	[]entity.User
// @Failure	404	{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router       /users/{id}/tasks [get]
func (h *Handler) getTasksByUser(c *gin.Context) {
	h.getTasksByUserId(c, "assignee", "id")
}

// @Summary get user by id
// @Tags users
// @Product json
// @Param		id           path		string	true	"path param"
// @Success	200		{object}	entity.User
// @Failure	404	{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router       /users/{id} [get]
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

// @Summary get user by name
// @Tags users
// @Product json
// @Param		name query		string	true	"path param"
// @Success	200		{object}	entity.User
// @Failure	404	{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router       /users/search/name [get]
func (h *Handler) getUserByName(c *gin.Context) {
	h.getUser(c, userName, "Name query parameter is required", h.service.User.GetUserByName)
}

// @Summary get user by email
// @Tags users
// @Product json
// @Param		email query		string	true	"path param"
// @Success	200		{object}	entity.User
// @Failure	404	{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router       /users/search/email [get]
func (h *Handler) getUserByEmail(c *gin.Context) {
	h.getUser(c, userEmail, "Email query parameter is required", h.service.User.GetUserByEmail)
}

// @Summary update user
// @Tags users
// @Produce json
// @Param id path string true "user id"
// @Param request body entity.User true "body param"
// @Success 200 {string} string "User updated successfully"
// @Failure 404 {object} response.Object
// @Failure 500 {object} response.Object
// @Router /users/{id} [put]
func (h *Handler) updateUser(c *gin.Context) {
	id := c.Param("id") // Corrected param key
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

// @Summary delete user
// @Tags users
// @Product json
// @Param		id path		string	true	"path param"
// @Success	200		{string}  string    "User deleted successfully"
// @Failure	404	{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router       /users/{id} [delete]
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

// @Summary get all users
// @Tags users
// @Product json
// @Success	200	{array}		entity.User
// @Failure	404	{object}	response.Object
// @Failure	500		{object}	response.Object
// @Router       /users [get]
func (h *Handler) getAllUsers(c *gin.Context) {
	lists, err := h.service.User.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lists)
}
