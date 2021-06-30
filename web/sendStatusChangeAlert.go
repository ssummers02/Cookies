package web

import (
	"log"
	"ssummers02/Cookies/bot"
	"ssummers02/Cookies/db"
	"strconv"
)

func sendChangeStatusAlert(taskID string) {
	var alert bot.StatusChangeAlert
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
	alert.Status = convertStatusIDtoStatusName(task.Status)
	bot.PostChangeStatus(alert)
}

func convertStatusIDtoStatusName(statusID uint) string {
	statusName := ""
	switch statusID {
	case 0:
		statusName = "Открыт"
	case 1:
		statusName = "Выполнен"
	case 2:
		statusName = "Требует уточнения"
	case 3:
		statusName = "Отклонен"
	case 4:
		statusName = "Отменен пользователем"
	}
	return statusName
}
