package handler

import (
	"github.com/Aytya/projects-manager-HL/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	user := router.Group("/users")
	{
		user.POST("/", h.createUser)                //создать нового пользователя.
		user.GET("/", h.getAllUsers)                //получить список всех пользователей.
		user.GET("/:id", h.getUserById)             //получить данные конкретного пользователя.
		user.PUT("/:id", h.updateUser)              //обновить данные конкретного пользователя.
		user.DELETE("/:id", h.deleteUser)           //удалить конкретного пользователя.
		user.GET("/:id/tasks", h.getTasksByUser)    //получить список задач конкретного пользователя.
		user.GET("/search/name", h.getUserByName)   //найти пользователей по имени.
		user.GET("/search/email", h.getUserByEmail) //найти пользователей по электронной почте.
	}

	task := router.Group("/tasks")
	{
		task.GET("/", h.getAllTasks)                                //получить список всех задач.
		task.POST("/", h.createTask)                                //создать новую задачу.
		task.GET("/:id", h.getTaskById)                             //получить данные конкретной задачи.
		task.PUT("/:id", h.updateTaskById)                          //обновить данные конкретной задачи.
		task.DELETE("/:id", h.deleteTaskById)                       //удалить конкретную задачу.
		task.GET("/search", h.getTaskByTitle)                       //найти задачи по названию.
		task.GET("/search/status", h.getTaskByStatus)               //найти задачи по состоянию.
		task.GET("/search/priority", h.getTaskByPriority)           //найти задачи по приоритету.
		task.GET("/search/:assignee", h.getTasksByAssignee)         // найти задачи по идентификатору ответственного.
		task.GET("/search/project/:project", h.getTasksByProjectId) //найти задачи по идентификатору проекта.
	}

	project := router.Group("/projects")
	{
		project.GET("/", h.getListOfProjects)              //получить список всех проектов.
		project.POST("/", h.createProject)                 //создать новый проект.
		project.GET("/:id", h.getProjectById)              //получить данные конкретного проекта.
		project.PUT("/:id", h.updateProject)               //обновить данные конкретного проекта.
		project.DELETE("/:id", h.deleteProject)            //удалить конкретный проект.
		project.GET("/:id/tasks", h.getProjectTasks)       //получить список задач в проекте.
		project.GET("/search/:title", h.getProjectByTitle) //найти проекты по названию.
		project.GET("/search", h.getProjectByManagerId)    //найти проекты по идентификатору менеджера.
	}

	return router
}
