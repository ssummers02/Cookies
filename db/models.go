package db

import "time"

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `json:"user"`
	Name      string `json:"name"`
	Floor     int    `json:"floor"`
	Room      string `json:"room"`
	Text      string `json:"text"`
	Status    uint   `json:"status"` // 1: New, 2: Done // Add if we have time for it)
	CreatedAt time.Time
}
type Users struct {
	UserID       int    `json:"user"`
	LastMessages string `json:"lastMessages"`
	Floor        int    `json:"floor"`
	Room         string `json:"room"`
}

type ArrayTask struct {
	Tasks []Task
}
