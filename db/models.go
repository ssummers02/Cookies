package db

import "time"

const (
	Created            uint = 0
	Completed          uint = 1
	NeedsClarification uint = 2
	Canceled           uint = 3
	CanceledByUser     uint = 4
)

type StatusChangeAlert struct {
	RecipientUserID int
	TaskID          string
	TaskText        string
	Status          string
}

type Task struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `json:"user"`
	Name      string `json:"name"`
	Floor     int    `json:"floor"`
	Room      string `json:"room"`
	Text      string `json:"text"`
	Status    uint   `json:"status"`
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
