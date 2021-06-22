package db

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

const limit = 100 // TODO: add .env

type Task struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `json:"user,string"`
	Room      uint      `json:"room,string"`
	Status    uint      `json:"status,string"` // 1: New, 2: Done // Add if we have time for it)
	CreatedAt time.Time `json:"time"`
}

func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	db.AutoMigrate(&Task{})
}

func CreateTask(task Task) error {
	return db.Create(&task).Error
}

func GetTask(id uint) Task {
	var task Task
	db.First(&task, id)
	return task
}

func GetAllTasks() []Task {
	var tasks []Task
	db.Limit(limit).Order("created_at, room").Find(&tasks)
	return tasks
}

func GetTaskInRoom(room uint) []Task {
	var tasks []Task
	db.Where(&Task{Room: room}, "room").Order("created_at, room").Find(&tasks)
	return tasks
}
