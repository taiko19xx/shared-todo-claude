package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shared-todo-backend/database"
	"shared-todo-backend/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *HandlerTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerTestSuite) SetupTest() {
	// テスト用データベースをセットアップ
	db, err := database.SetupTestDatabase()
	suite.Require().NoError(err)
	database.DB = db

	// Ginルーターをセットアップ
	suite.router = gin.New()
	suite.router.POST("/api/lists", CreateList)
	suite.router.GET("/api/lists/:listId/users/:userId", GetListData)
	suite.router.PUT("/api/lists/:listId/memo", UpdateListMemo)
	suite.router.POST("/api/lists/:listId/users", InviteUser)
	suite.router.PUT("/api/lists/:listId/users/:userId/name", UpdateUserName)
	suite.router.POST("/api/lists/:listId/todos", CreateTodo)
	suite.router.PUT("/api/todos/:todoId/status/:userId", UpdateTodoUserStatus)
}

func (suite *HandlerTestSuite) TearDownTest() {
	// テスト後のクリーンアップ
	if database.DB != nil {
		database.CleanupTestDatabase(database.DB)
	}
}

func (suite *HandlerTestSuite) TestCreateList() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/lists", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response, "listId")
	assert.Contains(suite.T(), response, "userId")
	assert.NotEmpty(suite.T(), response["listId"])
	assert.NotEmpty(suite.T(), response["userId"])
}

func (suite *HandlerTestSuite) TestGetListData() {
	// テストデータを作成
	list := models.List{ID: "test-list-id", Memo: "test memo"}
	database.DB.Create(&list)

	user := models.User{ID: "test-user-id", ListID: "test-list-id", DisplayName: "Test User"}
	database.DB.Create(&user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/lists/test-list-id/users/test-user-id", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response, "users")
	assert.Contains(suite.T(), response, "todos")
	assert.Contains(suite.T(), response, "memo")
	assert.Equal(suite.T(), "test memo", response["memo"])
}

func (suite *HandlerTestSuite) TestGetListDataNotFound() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/lists/nonexistent/users/nonexistent", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *HandlerTestSuite) TestUpdateListMemo() {
	// テストデータを作成
	list := models.List{ID: "test-list-id", Memo: ""}
	database.DB.Create(&list)

	payload := map[string]string{"memo": "Updated memo"}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/lists/test-list-id/memo", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated memo", response["memo"])

	// データベースで確認
	var updatedList models.List
	database.DB.First(&updatedList, "id = ?", "test-list-id")
	assert.Equal(suite.T(), "Updated memo", updatedList.Memo)
}

func (suite *HandlerTestSuite) TestCreateTodo() {
	// テストデータを作成
	list := models.List{ID: "test-list-id", Memo: ""}
	database.DB.Create(&list)

	user := models.User{ID: "test-user-id", ListID: "test-list-id", DisplayName: "Test User"}
	database.DB.Create(&user)

	payload := map[string]interface{}{
		"title":    "Test Todo",
		"priority": "high",
		"dueDate":  "2025-06-10",
	}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/lists/test-list-id/todos", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response models.Todo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Todo", response.Title)
	assert.Equal(suite.T(), "high", response.Priority)
	assert.Equal(suite.T(), "test-list-id", response.ListID)
	assert.False(suite.T(), response.IsCompleted)
}

func (suite *HandlerTestSuite) TestCreateTodoInvalidData() {
	// テストデータを作成
	list := models.List{ID: "test-list-id", Memo: ""}
	database.DB.Create(&list)

	// タイトルなしのToDo
	payload := map[string]interface{}{
		"priority": "high",
	}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/lists/test-list-id/todos", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *HandlerTestSuite) TestInviteUser() {
	// テストデータを作成
	list := models.List{ID: "test-list-id", Memo: ""}
	database.DB.Create(&list)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/lists/test-list-id/users", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Contains(suite.T(), response, "userId")
	assert.Contains(suite.T(), response, "url")
	assert.NotEmpty(suite.T(), response["userId"])
	assert.Contains(suite.T(), response["url"], response["userId"])
}

func (suite *HandlerTestSuite) TestUpdateUserName() {
	// テストデータを作成
	list := models.List{ID: "test-list-id", Memo: ""}
	database.DB.Create(&list)

	user := models.User{ID: "test-user-id", ListID: "test-list-id", DisplayName: ""}
	database.DB.Create(&user)

	payload := map[string]string{"name": "Updated Name"}
	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/lists/test-list-id/users/test-user-id/name", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Name", response["name"])
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}