package models

import (
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ModelsTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *ModelsTestSuite) SetupTest() {
	// インメモリSQLiteデータベースを使用
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// テーブル作成
	err = db.AutoMigrate(&List{}, &User{}, &Todo{}, &TodoUserStatus{})
	suite.Require().NoError(err)

	suite.db = db
}

func (suite *ModelsTestSuite) TearDownTest() {
	if suite.db != nil {
		// エラーを無視してテーブルをクリア
		suite.db.Exec("DELETE FROM todo_user_status")
		suite.db.Exec("DELETE FROM todos")
		suite.db.Exec("DELETE FROM users")
		suite.db.Exec("DELETE FROM lists")
	}
}

func (suite *ModelsTestSuite) TestListModel() {
	list := List{
		ID:   "test-list-id",
		Memo: "Test memo",
	}

	// Create
	err := suite.db.Create(&list).Error
	assert.NoError(suite.T(), err)

	// Read
	var retrievedList List
	err = suite.db.First(&retrievedList, "id = ?", "test-list-id").Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-list-id", retrievedList.ID)
	assert.Equal(suite.T(), "Test memo", retrievedList.Memo)
	assert.False(suite.T(), retrievedList.CreatedAt.IsZero())
	assert.False(suite.T(), retrievedList.UpdatedAt.IsZero())

	// Update
	retrievedList.Memo = "Updated memo"
	err = suite.db.Save(&retrievedList).Error
	assert.NoError(suite.T(), err)

	var updatedList List
	err = suite.db.First(&updatedList, "id = ?", "test-list-id").Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated memo", updatedList.Memo)
}

func (suite *ModelsTestSuite) TestUserModel() {
	// まずListを作成
	list := List{ID: "test-list-id", Memo: ""}
	err := suite.db.Create(&list).Error
	assert.NoError(suite.T(), err)

	user := User{
		ID:          "test-user-id",
		ListID:      "test-list-id",
		DisplayName: "Test User",
	}

	// Create
	err = suite.db.Create(&user).Error
	assert.NoError(suite.T(), err)

	// Read
	var retrievedUser User
	err = suite.db.First(&retrievedUser, "id = ?", "test-user-id").Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-user-id", retrievedUser.ID)
	assert.Equal(suite.T(), "test-list-id", retrievedUser.ListID)
	assert.Equal(suite.T(), "Test User", retrievedUser.DisplayName)
	assert.False(suite.T(), retrievedUser.CreatedAt.IsZero())
}

func (suite *ModelsTestSuite) TestTodoModel() {
	// まずListを作成
	list := List{ID: "test-list-id", Memo: ""}
	err := suite.db.Create(&list).Error
	assert.NoError(suite.T(), err)

	dueDate := time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC)
	todo := Todo{
		ListID:      "test-list-id",
		Title:       "Test Todo",
		Priority:    "high",
		DueDate:     &dueDate,
		IsCompleted: false,
	}

	// Create
	err = suite.db.Create(&todo).Error
	assert.NoError(suite.T(), err)
	assert.NotZero(suite.T(), todo.ID)

	// Read
	var retrievedTodo Todo
	err = suite.db.First(&retrievedTodo, todo.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-list-id", retrievedTodo.ListID)
	assert.Equal(suite.T(), "Test Todo", retrievedTodo.Title)
	assert.Equal(suite.T(), "high", retrievedTodo.Priority)
	assert.Equal(suite.T(), dueDate.Format("2006-01-02"), retrievedTodo.DueDate.Format("2006-01-02"))
	assert.False(suite.T(), retrievedTodo.IsCompleted)
	assert.False(suite.T(), retrievedTodo.CreatedAt.IsZero())
	assert.False(suite.T(), retrievedTodo.UpdatedAt.IsZero())

	// Update
	retrievedTodo.IsCompleted = true
	err = suite.db.Save(&retrievedTodo).Error
	assert.NoError(suite.T(), err)

	var updatedTodo Todo
	err = suite.db.First(&updatedTodo, todo.ID).Error
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), updatedTodo.IsCompleted)
}

func (suite *ModelsTestSuite) TestTodoUserStatusModel() {
	// 前提データの作成
	list := List{ID: "test-list-id", Memo: ""}
	err := suite.db.Create(&list).Error
	assert.NoError(suite.T(), err)

	user := User{ID: "test-user-id", ListID: "test-list-id", DisplayName: "Test User"}
	err = suite.db.Create(&user).Error
	assert.NoError(suite.T(), err)

	todo := Todo{ListID: "test-list-id", Title: "Test Todo", Priority: "medium", IsCompleted: false}
	err = suite.db.Create(&todo).Error
	assert.NoError(suite.T(), err)

	checkedAt := time.Now()
	status := TodoUserStatus{
		TodoID:    todo.ID,
		UserID:    "test-user-id",
		IsChecked: true,
		CheckedAt: &checkedAt,
	}

	// Create
	err = suite.db.Create(&status).Error
	assert.NoError(suite.T(), err)

	// Read
	var retrievedStatus TodoUserStatus
	err = suite.db.Where("todo_id = ? AND user_id = ?", todo.ID, "test-user-id").First(&retrievedStatus).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), todo.ID, retrievedStatus.TodoID)
	assert.Equal(suite.T(), "test-user-id", retrievedStatus.UserID)
	assert.True(suite.T(), retrievedStatus.IsChecked)
	assert.NotNil(suite.T(), retrievedStatus.CheckedAt)
}

func (suite *ModelsTestSuite) TestToDoModelWithNullDueDate() {
	// まずListを作成
	list := List{ID: "test-list-id", Memo: ""}
	err := suite.db.Create(&list).Error
	assert.NoError(suite.T(), err)

	todo := Todo{
		ListID:      "test-list-id",
		Title:       "Test Todo without due date",
		Priority:    "medium",
		DueDate:     nil, // 期限なし
		IsCompleted: false,
	}

	// Create
	err = suite.db.Create(&todo).Error
	assert.NoError(suite.T(), err)

	// Read
	var retrievedTodo Todo
	err = suite.db.First(&retrievedTodo, todo.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Todo without due date", retrievedTodo.Title)
	assert.Nil(suite.T(), retrievedTodo.DueDate)
}

func TestModelsTestSuite(t *testing.T) {
	suite.Run(t, new(ModelsTestSuite))
}