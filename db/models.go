package db

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type Task struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	Done      bool
	CreatedAt time.Time
	Deadline  time.Time
}

func InitDB() {
	var err error
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.AutoMigrate(&Task{})
}

func CreateOrder(task Task) error {
	return db.Create(&task).Error
}

func GetOrder(id uint) Task {
	var task Task
	db.First(&task, id)
	return task
}
