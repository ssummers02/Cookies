package db

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var limit int

func InitDB(dbName string, lim int) {
	var err error
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}
	limit = lim
	db.AutoMigrate(&Task{}, &Users{})
}

func CreateTask(task Task) error {
	task.CreatedAt = time.Now()
	return db.Create(&task).Error
}

func CreateUsers(user Users) error {
	return db.Create(&user).Error
}

func GetUsers(id int) (Users, error) {
	var user Users
	res := db.First(&user, id)
	return user, res.Error
}

func ChangeFloor(id int, n int) error {
	return db.Model(&Users{}).Where("user_id = ?", id).Update("floor", n).Error
}

func ChangeRoom(id int, n string) error {
	return db.Model(&Users{}).Where("user_id = ?", id).Update("room", n).Error
}

func ChangeMessage(id int, s string) error {
	return db.Model(&Users{}).Where("user_id = ?", id).Update("last_messages", s).Error
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
	res := db.Limit(countTasks).Where("user_id = ?", userId).Order("created_at desc, room").Find(&tasks)
	return tasks, res.Error
}

func GetTaskInRoom(room string) ([]Task, error) {
	var tasks []Task
	res := db.Limit(limit).Where("room = ?", room).Order("created_at, room").Find(&tasks)
	return tasks, res.Error
}

func DeleteTask(id uint) error {
	task, err := GetTask(id)
	if err != nil {
		return err
	}
	return db.Delete(&task).Error
}
