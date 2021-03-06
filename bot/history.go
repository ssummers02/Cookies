package bot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"ssummers02/Cookies/db"
)

func getHistory(peerId int) db.ArrayTask {
	var userHistory db.ArrayTask

	resp, err := http.Get("http://" + port + "/api/user/" + strconv.Itoa(peerId) + "/5")

	if err != nil {
		log.WithFields(log.Fields{
			"package": "history",
			"func":    "getHistory",
			"error":   err,
		}).Warning("error get history")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.WithFields(log.Fields{
			"package": "history",
			"func":    "getHistory",
			"error":   err,
		}).Warning("error get history")
	}
	err = json.Unmarshal(body, &userHistory)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "history",
			"func":    "getHistory",
			"error":   err,
		}).Warning("err json.Unmarshal")
	}

	return userHistory
}

func GetActiveHistory(PeerID int) db.ArrayTask {
	var userHistory db.ArrayTask

	resp, err := http.Get("http://" + port + "/api/useractive/" + strconv.Itoa(PeerID) + "/5")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "history",
			"func":    "GetActiveHistory",
			"error":   err,
		}).Warning("err ioutil.ReadAll")
	}

	err = json.Unmarshal(body, &userHistory)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "history",
			"func":    "GetActiveHistory",
			"error":   err,
		}).Warning("err json.Unmarshal")
	}

	return userHistory
}

func postHistoryForUser(vk *api.VK, peerId int) {
	userHistory := getHistory(peerId)

	if len(userHistory.Tasks) == 0 {
		postMessagesAndKeyboard(vk, peerId, "Заказов нет", getGeneralKeyboard(false))
		return
	}
	for i := 0; i < len(userHistory.Tasks); i++ {
		createMessage := "№" + strconv.Itoa(int(userHistory.Tasks[i].ID)) + ": " + userHistory.Tasks[i].Text + " - " + FindOutTheStatus(userHistory.Tasks[i].Status) + "\n"
		postMessagesAndKeyboard(vk, peerId, createMessage, getGeneralKeyboard(false))
	}
}
func selectDeleteHistory(vk *api.VK, peerId int) {
	userHistory := GetActiveHistory(peerId)

	if len(userHistory.Tasks) == 0 {
		postMessagesAndKeyboard(vk, peerId, "Заказов нет", getGeneralKeyboard(false))
		return
	}
	k := object.NewMessagesKeyboardInline()
	k.AddRow()

	for i := 0; i < len(userHistory.Tasks); i++ {
		id := strconv.Itoa(int(userHistory.Tasks[i].ID))
		k.AddTextButton(id, ``, `primary`)
		createMessage := "№" + id + ": " + userHistory.Tasks[i].Text + " - " + FindOutTheStatus(userHistory.Tasks[i].Status) + "\n"
		postAndSendMessages(vk, peerId, createMessage)
	}
	k.AddRow()

	k.AddTextButton(`Вернуться назад`, ``, `positive`)

	postMessagesAndKeyboard(vk, peerId, "Выбери заказ который хочешь отменить", k)
}
