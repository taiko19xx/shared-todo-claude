package database

import (
	"shared-todo-backend/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetupTestDatabase() (*gorm.DB, error) {
	// インメモリSQLiteデータベースを使用
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// テーブル作成
	err = db.AutoMigrate(
		&models.List{},
		&models.User{},
		&models.Todo{},
		&models.TodoUserStatus{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CleanupTestDatabase(db *gorm.DB) error {
	// エラーを無視してテーブルをクリア
	db.Exec("DELETE FROM todo_user_status")
	db.Exec("DELETE FROM todos")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM lists")
	return nil
}