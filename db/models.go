package db

import "time"

const (
	Created            uint = 0
	Completed          uint = 1
	NeedsClarification uint = 2
	Canceled           uint = 3
	CanceledByUser     uint = 4
)

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
