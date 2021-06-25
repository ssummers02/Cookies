package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"ssummers02/Cookies/db"
)

func GetHistory(port string, PeerID int) db.ArrayTask {
	var userHistory db.ArrayTask

	resp, err := http.Get("http://" + port + "/api/user/" + strconv.Itoa(PeerID) + "/5")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &userHistory)

	return userHistory
}

func PostHistoryForUser(vk *api.VK, PeerID int) {
	port := os.Getenv("ADDRESS")
	userHistory := GetHistory(port, PeerID)

	if len(userHistory.Tasks) == 0 {
		PostMessagesAndKeyboard(vk, PeerID, "Заказов нет", GetGeneralKeyboard(false))
		return
	}
	for i := 0; i < len(userHistory.Tasks); i++ {
		createMessage := "№" + strconv.Itoa(int(userHistory.Tasks[i].ID)) + ": " + userHistory.Tasks[i].Text + " - " + findOutTheStatus(userHistory.Tasks[i].Status) + "\n"
		PostMessagesAndKeyboard(vk, PeerID, createMessage, GetGeneralKeyboard(false))
	}
}
func SelectDeleteHistory(vk *api.VK, PeerID int) {
	port := os.Getenv("ADDRESS")
	userHistory := GetHistory(port, PeerID)

	if len(userHistory.Tasks) == 0 {
		PostMessagesAndKeyboard(vk, PeerID, "Заказов нет", GetGeneralKeyboard(false))
		return
	}
	k := object.NewMessagesKeyboardInline()
	k.AddRow()

	for i := 0; i < len(userHistory.Tasks); i++ {
		id := strconv.Itoa(int(userHistory.Tasks[i].ID))
		k.AddTextButton(id, ``, `primary`)
		createMessage := "№" + id + ": " + userHistory.Tasks[i].Text + " - " + findOutTheStatus(userHistory.Tasks[i].Status) + "\n"
		PostAndSendMessages(vk, PeerID, createMessage)
	}
	k.AddRow()

	k.AddTextButton(`Вернуться назад`, ``, `positive`)

	PostMessagesAndKeyboard(vk, PeerID, "Выбери заказ который хочешь отменить", k)
}
