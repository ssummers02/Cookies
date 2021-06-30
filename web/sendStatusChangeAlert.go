package web

import (
	"log"
	"strconv"

	"ssummers02/Cookies/bot"
	"ssummers02/Cookies/db"
)

func sendChangeStatusAlert(taskID string) {
	var alert db.StatusChangeAlert
	var task db.Task
	var err error
	var taskIDInt int
	taskIDInt, err = strconv.Atoi(taskID)
	if err != nil {
		log.Print(err)
		return
	}
	task, err = db.GetTask(uint(taskIDInt))
	if err != nil {
		log.Print(err)
		return
	}
	alert.RecipientUserID = int(task.UserID)
	alert.TaskID = strconv.Itoa(int(task.ID))
	alert.TaskText = task.Text
	alert.Status = bot.FindOutTheStatus(task.Status)
	bot.PostChangeStatus(alert)
}
