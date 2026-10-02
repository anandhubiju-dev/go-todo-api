package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"

	"todo-api/config"
	"todo-api/middleware"
	"todo-api/model"
	"todo-api/repository"
	"todo-api/service"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		fmt.Println("Configuration error:", err)
		return
	}
	client, todoCollection, err := connectDB(cfg)

	_ = client // To avoid unused variable error

	if err != nil {
		fmt.Println("Failed to connect to MongoDB:", err)
		return
	}

	fmt.Println("Connected to MongoDB:")
	fmt.Println(todoCollection.Name())

	todoRepo := repository.NewTodoRepository(todoCollection)
	todoService := service.NewTodoService(todoRepo)

	router := gin.Default()

	router.Use(middleware.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.GET("/todos/:id", func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := bson.ObjectIDFromHex(idStr)

		if err != nil {
			c.JSON(400, gin.H{
				"error": "invalid todo ID",
			})
			return
		}

		todo, err := todoService.GetByID(
			c.Request.Context(),
			id,
		)

		if err != nil {
			if errors.Is(err, service.ErrTodoNotFound) {
				c.JSON(404, gin.H{
					"error": "todo not found",
				})
				return
			}

			c.JSON(500, gin.H{
				"error": "failed to fetch todo",
			})
			return
		}

		c.JSON(200, todo)
	})

	router.GET("/todos", func(c *gin.Context) {
		filter := bson.M{}

		completedStr, exists := c.GetQuery("completed")

		if exists {
			completed, err := strconv.ParseBool(completedStr)

			if err != nil {
				c.JSON(400, gin.H{
					"error": "invalid 'complete' parameter",
				})
				return
			}

			filter["completed"] = completed
		}

		todos, err := todoService.GetAll(
			c.Request.Context(),
			filter,
		)

		if err != nil {
			c.JSON(500, gin.H{
				"error": "faoled to fetch todos",
			})
			return
		}
		c.JSON(200, todos)
	})

	router.PUT("/todos/:id", func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := bson.ObjectIDFromHex(idStr)

		if err != nil {
			c.JSON(400, gin.H{
				"error": "invalid todo ID",
			})
			return
		}

		var updatedTodo model.Todo

		err = c.ShouldBindJSON(&updatedTodo)

		if err != nil {
			c.JSON(400, gin.H{
				"error": "invalid request body",
			})
			return
		}

		err = todoService.Update(
			c.Request.Context(),
			id,
			updatedTodo,
		)

		if err != nil {
			if errors.Is(err, service.ErrTodoNotFound) {
				c.JSON(404, gin.H{
					"error": "todo not found",
				})
				return
			}
			c.JSON(500, gin.H{
				"error": "failed to update todo",
			})
			return
		}

		updatedTodo.ID = id

		c.JSON(200, updatedTodo)
	})

	router.DELETE("/todos/:id", func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := bson.ObjectIDFromHex(idStr)

		if err != nil {
			c.JSON(400, gin.H{
				"error": "invalid todo ID",
			})
			return
		}

		err = todoService.Delete(
			c.Request.Context(),
			id,
		)

		if err != nil {
			if errors.Is(err, service.ErrTodoNotFound) {
				c.JSON(404, gin.H{
					"error": "todo not found",
				})
				return
			}

			c.JSON(500, gin.H{
				"error": "failed to delete todo",
			})
			return
		}

		c.JSON(200, gin.H{
			"message": "todo deleted",
		})
	})

	router.POST("/todos", func(c *gin.Context) {
		var todo model.Todo

		err := c.ShouldBindJSON(&todo)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid JSON"})
			return
		}

		insertedID, err := todoService.Create(
			c.Request.Context(),
			todo,
		)

		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create todo"})
			return
		}

		todo.ID = insertedID

		c.JSON(201, todo)
	})

	router.Run(":" + cfg.Port)
}
