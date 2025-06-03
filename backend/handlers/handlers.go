package handlers

import (
	"net/http"
	"strconv"
	"time"
	"shared-todo-backend/database"
	"shared-todo-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateList creates a new list and user
func CreateList(c *gin.Context) {
	listID := uuid.New().String()
	userID := uuid.New().String()

	list := models.List{
		ID:   listID,
		Memo: "",
	}

	user := models.User{
		ID:          userID,
		ListID:      listID,
		DisplayName: "",
	}

	if err := database.DB.Create(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create list"})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"listId": listID,
		"userId": userID,
	})
}

// GetListData gets list information and user information
func GetListData(c *gin.Context) {
	listID := c.Param("listId")
	userID := c.Param("userId")

	// Check if user exists in the list
	var user models.User
	if err := database.DB.Where("id = ? AND list_id = ?", userID, listID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in this list"})
		return
	}

	// Get list with memo
	var list models.List
	if err := database.DB.First(&list, "id = ?", listID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	// Get all users in the list
	var users []models.User
	database.DB.Where("list_id = ?", listID).Find(&users)

	// Get all todos with user statuses
	var todos []models.Todo
	database.DB.Where("list_id = ?", listID).Preload("UserStatuses").Find(&todos)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"todos": todos,
		"memo":  list.Memo,
	})
}

// UpdateListMemo updates the memo of a list
func UpdateListMemo(c *gin.Context) {
	listID := c.Param("listId")

	// Check if list exists
	var list models.List
	if err := database.DB.First(&list, "id = ?", listID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	var req struct {
		Memo string `json:"memo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate memo length
	if len(req.Memo) > 5000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Memo must be 5000 characters or less"})
		return
	}

	if err := database.DB.Model(&models.List{}).Where("id = ?", listID).Update("memo", req.Memo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update memo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"memo": req.Memo})
}

// InviteUser creates a new user for the list
func InviteUser(c *gin.Context) {
	listID := c.Param("listId")

	// Check if list exists
	var list models.List
	if err := database.DB.First(&list, "id = ?", listID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	userID := uuid.New().String()
	user := models.User{
		ID:          userID,
		ListID:      listID,
		DisplayName: "",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Create todo user status records for existing todos
	var todos []models.Todo
	database.DB.Where("list_id = ?", listID).Find(&todos)

	for _, todo := range todos {
		status := models.TodoUserStatus{
			TodoID:    todo.ID,
			UserID:    userID,
			IsChecked: false,
		}
		database.DB.Create(&status)
	}

	c.JSON(http.StatusCreated, gin.H{
		"userId": userID,
		"url":    "/" + listID + "/" + userID,
	})
}

// UpdateUserName updates user display name
func UpdateUserName(c *gin.Context) {
	listID := c.Param("listId")
	userID := c.Param("userId")

	// Check if user exists in the list
	var user models.User
	if err := database.DB.Where("id = ? AND list_id = ?", userID, listID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found in this list"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate name length
	if len(req.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Display name must be 100 characters or less"})
		return
	}

	if err := database.DB.Model(&user).Update("display_name", req.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user name"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"name": req.Name})
}

// CreateTodo creates a new todo
func CreateTodo(c *gin.Context) {
	listID := c.Param("listId")

	// Check if list exists
	var list models.List
	if err := database.DB.First(&list, "id = ?", listID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "List not found"})
		return
	}

	var req struct {
		Title    string  `json:"title" binding:"required"`
		Priority string  `json:"priority"`
		DueDate  *string `json:"dueDate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Validate title
	if len(req.Title) > 255 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title must be 255 characters or less"})
		return
	}

	// Validate priority
	if req.Priority == "" {
		req.Priority = "medium"
	}
	if req.Priority != "high" && req.Priority != "medium" && req.Priority != "low" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Priority must be 'high', 'medium', or 'low'"})
		return
	}

	// Parse due date
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		dueDate = &parsedDate
	}

	todo := models.Todo{
		ListID:      listID,
		Title:       req.Title,
		Priority:    req.Priority,
		DueDate:     dueDate,
		IsCompleted: false,
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
		return
	}

	// Create todo user status records for all users in the list
	var users []models.User
	database.DB.Where("list_id = ?", listID).Find(&users)

	for _, user := range users {
		status := models.TodoUserStatus{
			TodoID:    todo.ID,
			UserID:    user.ID,
			IsChecked: false,
		}
		database.DB.Create(&status)
	}

	c.JSON(http.StatusCreated, todo)
}

// UpdateTodoUserStatus updates user's check status for a todo
func UpdateTodoUserStatus(c *gin.Context) {
	todoIDStr := c.Param("todoId")
	userID := c.Param("userId")

	todoID, err := strconv.ParseUint(todoIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID format"})
		return
	}

	// Check if todo exists
	var todo models.Todo
	if err := database.DB.First(&todo, todoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	// Check if user exists in the same list as the todo
	var user models.User
	if err := database.DB.Where("id = ? AND list_id = ?", userID, todo.ListID).First(&user).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not authorized to update this todo"})
		return
	}

	var req struct {
		Checked bool `json:"checked"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Update user status
	var status models.TodoUserStatus
	result := database.DB.Where("todo_id = ? AND user_id = ?", todoID, userID).First(&status)
	
	if result.Error != nil {
		// Create new status if not exists
		status = models.TodoUserStatus{
			TodoID:    uint(todoID),
			UserID:    userID,
			IsChecked: req.Checked,
		}
		if req.Checked {
			now := time.Now()
			status.CheckedAt = &now
		}
		database.DB.Create(&status)
	} else {
		// Update existing status
		status.IsChecked = req.Checked
		if req.Checked {
			now := time.Now()
			status.CheckedAt = &now
		} else {
			status.CheckedAt = nil
		}
		database.DB.Save(&status)
	}

	// Check if all users have checked this todo
	database.DB.First(&todo, todoID)

	var users []models.User
	database.DB.Where("list_id = ?", todo.ListID).Find(&users)

	var checkedCount int64
	database.DB.Model(&models.TodoUserStatus{}).Where("todo_id = ? AND is_checked = ?", todoID, true).Count(&checkedCount)

	isCompleted := int(checkedCount) == len(users)
	database.DB.Model(&todo).Update("is_completed", isCompleted)

	c.JSON(http.StatusOK, gin.H{"checked": req.Checked})
}