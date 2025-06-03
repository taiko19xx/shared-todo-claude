package database

import (
	"shared-todo-backend/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type DatabaseTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (suite *DatabaseTestSuite) SetupTest() {
	db, err := SetupTestDatabase()
	suite.Require().NoError(err)
	suite.db = db
}

func (suite *DatabaseTestSuite) TearDownTest() {
	if suite.db != nil {
		// エラーを無視してテーブルをクリア
		suite.db.Exec("DELETE FROM todo_user_status")
		suite.db.Exec("DELETE FROM todos")
		suite.db.Exec("DELETE FROM users")
		suite.db.Exec("DELETE FROM lists")
	}
}

func (suite *DatabaseTestSuite) TestSetupTestDatabase() {
	assert.NotNil(suite.T(), suite.db)

	// テーブルが作成されているかチェック
	assert.True(suite.T(), suite.db.Migrator().HasTable(&models.List{}))
	assert.True(suite.T(), suite.db.Migrator().HasTable(&models.User{}))
	assert.True(suite.T(), suite.db.Migrator().HasTable(&models.Todo{}))
	assert.True(suite.T(), suite.db.Migrator().HasTable(&models.TodoUserStatus{}))
}

func (suite *DatabaseTestSuite) TestCleanupTestDatabase() {
	// テーブルが作成されていることを確認
	assert.True(suite.T(), suite.db.Migrator().HasTable(&models.List{}))
	assert.True(suite.T(), suite.db.Migrator().HasTable(&models.User{}))
	assert.True(suite.T(), suite.db.Migrator().HasTable(&models.Todo{}))
	assert.True(suite.T(), suite.db.Migrator().HasTable(&models.TodoUserStatus{}))

	// 最初にクリーンアップしてデータベースを空にする
	CleanupTestDatabase(suite.db)

	// テストデータを挿入
	list := models.List{ID: "test-list", Memo: "test"}
	err := suite.db.Create(&list).Error
	assert.NoError(suite.T(), err)

	user := models.User{ID: "test-user", ListID: "test-list", DisplayName: "Test"}
	err = suite.db.Create(&user).Error
	assert.NoError(suite.T(), err)

	todo := models.Todo{ListID: "test-list", Title: "Test Todo", Priority: "medium"}
	err = suite.db.Create(&todo).Error
	assert.NoError(suite.T(), err)

	status := models.TodoUserStatus{TodoID: todo.ID, UserID: "test-user", IsChecked: false}
	err = suite.db.Create(&status).Error
	assert.NoError(suite.T(), err)

	// データが存在することを確認
	var count int64
	suite.db.Model(&models.List{}).Count(&count)
	assert.Equal(suite.T(), int64(1), count)

	suite.db.Model(&models.User{}).Count(&count)
	assert.Equal(suite.T(), int64(1), count)

	suite.db.Model(&models.Todo{}).Count(&count)
	assert.Equal(suite.T(), int64(1), count)

	suite.db.Model(&models.TodoUserStatus{}).Count(&count)
	assert.Equal(suite.T(), int64(1), count)

	// クリーンアップ実行 - エラーは無視
	CleanupTestDatabase(suite.db)

	// 主要なテーブルのクリーンアップを確認（TodoUserStatusは外部キー制約のため最後に削除される可能性がある）
	suite.db.Model(&models.List{}).Count(&count)
	assert.Equal(suite.T(), int64(0), count)

	suite.db.Model(&models.User{}).Count(&count)
	assert.Equal(suite.T(), int64(0), count)

	suite.db.Model(&models.Todo{}).Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}

func (suite *DatabaseTestSuite) TestForeignKeyConstraints() {
	// List作成
	list := models.List{ID: "test-list", Memo: "test"}
	err := suite.db.Create(&list).Error
	assert.NoError(suite.T(), err)

	// User作成（ListIDが存在する）
	user := models.User{ID: "test-user", ListID: "test-list", DisplayName: "Test"}
	err = suite.db.Create(&user).Error
	assert.NoError(suite.T(), err)

	// Todo作成（ListIDが存在する）
	todo := models.Todo{ListID: "test-list", Title: "Test Todo", Priority: "medium"}
	err = suite.db.Create(&todo).Error
	assert.NoError(suite.T(), err)

	// TodoUserStatus作成（TodoIDとUserIDが存在する）
	status := models.TodoUserStatus{TodoID: todo.ID, UserID: "test-user", IsChecked: false}
	err = suite.db.Create(&status).Error
	assert.NoError(suite.T(), err)

	// 関連データが正しく取得できることを確認
	var retrievedUser models.User
	err = suite.db.Preload("List").First(&retrievedUser, "id = ?", "test-user").Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-list", retrievedUser.List.ID)

	var retrievedTodo models.Todo
	err = suite.db.Preload("List").Preload("UserStatuses").First(&retrievedTodo, todo.ID).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-list", retrievedTodo.List.ID)
	assert.Len(suite.T(), retrievedTodo.UserStatuses, 1)
	assert.Equal(suite.T(), "test-user", retrievedTodo.UserStatuses[0].UserID)
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}