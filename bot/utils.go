package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"ssummers02/Cookies/db"
)

const (
	Created            uint = 0
	completed          uint = 1
	NeedsClarification uint = 2
	canceled           uint = 3
	canceledByUser     uint = 4
)

func PostChangeStatus(userid int, taskid string, status string) {
	vkKey := os.Getenv("VK_KEY")

	vk := api.NewVK(vkKey)
	PostAndSendMessages(vk, userid, "Заказ: "+taskid+"-"+status) // второй аргумент кому отдать изменение статуса
}

func PostNewTask(vk *api.VK, Message string, PeerID int, room string, floor int) {
	emp := &db.Task{UserID: uint(PeerID), Name: GetName(vk, PeerID), Room: room, Text: Message, Floor: floor}
	jsonData, _ := json.Marshal(emp)

	_, err := http.Post("http://"+port+"/api/add_task", "application/json",
		bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatal(err)
	}
}

func GetName(vk *api.VK, PeerID int) string {
	b := params.NewUsersGetBuilder()
	var id = []string{strconv.Itoa(PeerID)}
	b.UserIDs(id)

	resp, err := vk.UsersGet(b.Params)

	if err != nil {
		log.Fatal(err)
	}
	return resp[0].FirstName + " " + resp[0].LastName
}

func FindOutTheStatus(n uint) string {
	switch n {
	case 0:
		return "создана"
	case 1:
		return "выполнена"
	case 2:
		return "требует уточнения"
	case 3:
		return "отклонена"
	case 4:
		return "отменена пользователем"
	}
	return ""
}

func PostFloor(vk *api.VK, Message string, PeerID int) {
	floor, err := strconv.Atoi(Message)
	if err != nil { // если возникла ошибка
		PostAndSendMessages(vk, PeerID, "Неверный этаж, попробуй еще раз")
	} else {
		PostMessagesAndKeyboard(vk, PeerID, "Твой этаж:"+strconv.Itoa(floor)+"\nЧем я могу тебе помочь?", GetGeneralKeyboard(true))
		db.ChangeFloor(PeerID, floor)
	}
}
func ChangeStatus(vk *api.VK, Message string, PeerID int) string {
	userHistory := GetActiveHistory(PeerID)

	for _, task := range userHistory.Tasks {
		if strconv.Itoa(int(task.ID)) == Message {
			req, err := http.NewRequest(http.MethodPut, "http://"+port+"/api/task/status/"+Message+"/4", nil)
			if err != nil {
				fmt.Println(err)
			}
			_, err = http.DefaultClient.Do(req)
			PostMessagesAndKeyboard(vk, PeerID, "Твой заказ отменен", GetGeneralKeyboard(false))
			return Message
		}
	}
	PostMessagesAndKeyboard(vk, PeerID, "Этот заказ не может быть отменен", GetGeneralKeyboard(false))

	return Message
}
