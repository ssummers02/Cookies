package db

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var limit int

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `json:"user"`
	Room      uint   `json:"room"`
	Text      string `json:"text"`
	Status    uint   `json:"status"` // 1: New, 2: Done // Add if we have time for it)
	CreatedAt time.Time
}

func InitDB(dbName string, lim int) {
	var err error
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	limit = lim
	db.AutoMigrate(&Task{})
}

func CreateTask(task Task) error {
	task.CreatedAt = time.Now()
	return db.Create(&task).Error
}

func UpdateTask(task Task) error {
	return db.Save(&task).Error
}

func ChangeStatus(taskId string, status string) error {
	return db.Model(&Task{}).Where("id = ?", taskId).Update("status", status).Error
}

func GetTask(id uint) (Task, error) {
	var task Task
	res := db.First(&task, id)
	return task, res.Error
}

func GetAllTasks() ([]Task, error) {
	var tasks []Task
	res := db.Limit(limit).Order("created_at, room").Find(&tasks)
	return tasks, res.Error
}

func GetUserHistory(userId uint, countTasks int) ([]Task, error) {
	var tasks []Task
	res := db.Limit(countTasks).Where("user_id = ?", userId).Order("created_at, room").Find(&tasks)
	return tasks, res.Error
}

func GetTaskInRoom(room string) ([]Task, error) {
	var tasks []Task
	res := db.Limit(limit).Where("room = ?", room).Order("created_at, room").Find(&tasks)
	return tasks, res.Error
}

func DeleteTask(id uint) error {
	task, err := GetTask(uint(id))
	if err != nil {
		return err
	}
	return db.Delete(&task).Error
}
