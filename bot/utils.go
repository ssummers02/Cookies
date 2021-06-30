package bot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"ssummers02/Cookies/db"
)

func PostChangeStatus(alert db.StatusChangeAlert) {
	vkKey := os.Getenv("VK_KEY")
	vk := api.NewVK(vkKey)
	postAndSendMessages(vk, alert.RecipientUserID, "У заказа №"+alert.TaskID+"\n("+alert.TaskText+
		")\nизменился статус на "+alert.Status) // второй аргумент кому отдать изменение статуса
}

func postNewTask(vk *api.VK, message string, peerId int, room string, floor int) {
	emp := &db.Task{UserID: uint(peerId), Name: getName(vk, peerId), Room: room, Text: message, Floor: floor}
	jsonData, _ := json.Marshal(emp)

	_, err := http.Post("http://"+port+"/api/add_task", "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		log.WithFields(log.Fields{
			"package": "utils",
			"func":    "postNewTask",
			"error":   err,
		}).Warning("error post new task")
	}
}

func getName(vk *api.VK, peerId int) string {
	b := params.NewUsersGetBuilder()
	var id = []string{strconv.Itoa(peerId)}
	b.UserIDs(id)

	resp, err := vk.UsersGet(b.Params)

	if err != nil {
		log.WithFields(log.Fields{
			"package": "utils",
			"func":    "getName",
			"error":   err,
		}).Warning("get name")
	}
	return resp[0].FirstName + " " + resp[0].LastName
}

func FindOutTheStatus(n uint) string {
	switch n {
	case 0:
		return "создан"
	case 1:
		return "выполнен"
	case 2:
		return "требует уточнения"
	case 3:
		return "отклонён"
	case 4:
		return "отменён пользователем"
	}
	return ""
}

func postFloor(vk *api.VK, message string, peerId int) string {
	floor, err := strconv.Atoi(message)
	if err != nil { // если возникла ошибка
		postAndSendMessages(vk, peerId, "Неверный этаж, попробуй еще раз")
		return "Этаж"
	} else {
		err = db.ChangeFloor(peerId, floor)
		if err != nil {
			log.Print(err)
			postMessagesAndKeyboard(vk, peerId, "Произошла ошибка, попробуйте повторить позже.\nЧем я могу тебе помочь?", getGeneralKeyboard(true))
		}
		postMessagesAndKeyboard(vk, peerId, "Твой этаж: "+strconv.Itoa(floor)+"\nЧем я могу тебе помочь?", getGeneralKeyboard(true))
	}
	return ""
}
func changeStatus(vk *api.VK, message string, peerId int) string {
	userHistory := GetActiveHistory(peerId)
	for _, task := range userHistory.Tasks {
		if strconv.Itoa(int(task.ID)) == message {
			req, err := http.NewRequest(http.MethodPut, "http://"+port+"/api/task/status/"+message+"/4", nil)
			if err != nil {
				log.WithFields(log.Fields{
					"package": "utils",
					"func":    "changeStatus",
					"error":   err,
				}).Warning("change Status")
			}
			_, err = http.DefaultClient.Do(req)
			postMessagesAndKeyboard(vk, peerId, "Твой заказ отменен", getGeneralKeyboard(false))
			return message
		}
	}
	postMessagesAndKeyboard(vk, peerId, "Этот заказ не может быть отменен", getGeneralKeyboard(false))

	return message
}
