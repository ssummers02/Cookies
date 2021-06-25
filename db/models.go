package db

import (
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

var limit int

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `json:"user"`
	Floor     int    `json:"floor"`
	Room      string `json:"room"`
	Text      string `json:"text"`
	Status    uint   `json:"status"` // 1: New, 2: Done // Add if we have time for it)
	CreatedAt time.Time
}
type Users struct {
	UserID       int    `json:"user"`
	Name         string `json:"name"`
	LastMessages string `json:"LastMessages"`
	Floor        int    `json:"floor"`
	Room         string `json:"Room"`
}

type ArrayTask struct {
	Tasks []Task
}
