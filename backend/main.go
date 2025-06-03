package main

import (
	"log"
	"os"
	"shared-todo-backend/database"
	"shared-todo-backend/handlers"
	"shared-todo-backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// データベース初期化
	database.InitDatabase()

	// Ginエンジン初期化
	r := gin.Default()

	// CORS設定
	r.Use(middleware.SetupCORS())

	// APIルート
	api := r.Group("/api")
	{
		// リスト関連
		api.POST("/lists", handlers.CreateList)
		api.GET("/lists/:listId/users/:userId", handlers.GetListData)
		api.PUT("/lists/:listId/memo", handlers.UpdateListMemo)

		// ユーザー関連
		api.POST("/lists/:listId/users", handlers.InviteUser)
		api.PUT("/lists/:listId/users/:userId/name", handlers.UpdateUserName)

		// ToDo関連
		api.POST("/lists/:listId/todos", handlers.CreateTodo)
		api.PUT("/todos/:todoId/status/:userId", handlers.UpdateTodoUserStatus)
	}

	// ポート設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}