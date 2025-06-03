package models

import (
	"time"
)

type List struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Memo      string    `json:"memo" gorm:"default:''"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Users     []User    `json:"users,omitempty" gorm:"foreignKey:ListID"`
	Todos     []Todo    `json:"todos,omitempty" gorm:"foreignKey:ListID"`
}

type User struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	ListID      string    `json:"listId" gorm:"not null"`
	DisplayName string    `json:"displayName" gorm:"default:''"`
	CreatedAt   time.Time `json:"createdAt"`
	List        List      `json:"-" gorm:"foreignKey:ListID"`
}

type Todo struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	ListID       string            `json:"listId" gorm:"not null"`
	Title        string            `json:"title" gorm:"not null"`
	Priority     string            `json:"priority" gorm:"default:'medium';check:priority IN ('high', 'medium', 'low')"`
	DueDate      *time.Time        `json:"dueDate"`
	IsCompleted  bool              `json:"isCompleted" gorm:"default:false"`
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
	List         List              `json:"-" gorm:"foreignKey:ListID"`
	UserStatuses []TodoUserStatus  `json:"userStatuses,omitempty" gorm:"foreignKey:TodoID"`
}

type TodoUserStatus struct {
	TodoID    uint       `json:"todoId" gorm:"primaryKey"`
	UserID    string     `json:"userId" gorm:"primaryKey"`
	IsChecked bool       `json:"isChecked" gorm:"default:false"`
	CheckedAt *time.Time `json:"checkedAt"`
	Todo      Todo       `json:"-" gorm:"foreignKey:TodoID"`
	User      User       `json:"-" gorm:"foreignKey:UserID"`
}