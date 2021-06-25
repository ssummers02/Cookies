package bot

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"ssummers02/Cookies/db"
)

func PostNewTask(vk *api.VK, Message string, PeerID int, room string, floor int) {
	port := os.Getenv("ADDRESS")
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

func findOutTheStatus(n uint) string {
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
