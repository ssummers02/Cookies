package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"ssummers02/Cookies/db"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
)

// Отправка сообщения пользователю
func createAndSendMessages(vk *api.VK, PeerID int, text string) {
	rand.Seed(time.Now().UnixNano())
	b := params.NewMessagesSendBuilder()

	b.Message(text)
	b.RandomID(rand.Intn(2147483647))
	b.PeerID(PeerID)
	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.Fatal(err)
	}
}

// Отправка сообщения пользователю
func createAndSendMessagesAndKeyboard(vk *api.VK, PeerID int, text string, k *object.MessagesKeyboard) {
	rand.Seed(time.Now().UnixNano())
	b := params.NewMessagesSendBuilder()

	b.Keyboard(k)
	b.Message(text)
	b.RandomID(rand.Intn(2147483647))
	b.PeerID(PeerID)
	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.Fatal(err)
	}
}

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

func sendHistory(vk *api.VK, port string, PeerID int) {
	userHistory := GetHistory(port, PeerID)

	if len(userHistory.Tasks) == 0 {
		createAndSendMessagesAndKeyboard(vk, PeerID, "Заказов нет", createGeneralKeyboard(false))
		return
	}
	for i := 0; i < len(userHistory.Tasks); i++ {
		createMessage := "№" + strconv.Itoa(int(userHistory.Tasks[i].ID)) + ": " + userHistory.Tasks[i].Text + " - " + findOutTheStatus(userHistory.Tasks[i].Status) + "\n"
		createAndSendMessagesAndKeyboard(vk, PeerID, createMessage, createGeneralKeyboard(false))
	}
}
func SelectDeleteHistory(vk *api.VK, port string, PeerID int) {
	userHistory := GetHistory(port, PeerID)

	if len(userHistory.Tasks) == 0 {
		createAndSendMessagesAndKeyboard(vk, PeerID, "Заказов нет", createGeneralKeyboard(false))
		return
	}
	k := object.NewMessagesKeyboardInline()
	k.AddRow()

	for i := 0; i < len(userHistory.Tasks); i++ {
		id := strconv.Itoa(int(userHistory.Tasks[i].ID))
		k.AddTextButton(id, ``, `primary`)
		createMessage := "№" + id + ": " + userHistory.Tasks[i].Text + " - " + findOutTheStatus(userHistory.Tasks[i].Status) + "\n"
		createAndSendMessages(vk, PeerID, createMessage)
	}
	k.AddRow()

	k.AddTextButton(`Вернуться назад`, ``, `positive`)

	createAndSendMessagesAndKeyboard(vk, PeerID, "Выбери заказ который хочешь отменить", k)
}

func messageHandling(vk *api.VK, Message string, PeerID int) string {
	userStatus, _ := db.GetUsers(PeerID)

	port := os.Getenv("ADDRESS")

	if userStatus.Room == "0" && userStatus.LastMessages == "Кабинет" {
		db.ChangeRoom(PeerID, Message)
		createAndSendMessages(vk, PeerID, "Твой новый кабинет: "+Message+"\n Укажи этаж")
		return "Этаж"
	}

	if userStatus.Room == "0" && userStatus.LastMessages != "Кабинет" {
		createAndSendMessages(vk, PeerID, "Я тебя не знаю, давай познакомимься поближе\nУкажи номер своего кабинета")
		return "Кабинет"
	}
	if userStatus.LastMessages == "Этаж" {
		floor, err := strconv.Atoi(Message)
		if err != nil { // если возникла ошибка
			createAndSendMessages(vk, PeerID, "Неверный этаж, попробуй еще раз")
		} else {
			createAndSendMessagesAndKeyboard(vk, PeerID, "Твой этаж:"+strconv.Itoa(floor)+"\nЧем я могу тебе помочь?", createGeneralKeyboard(true))
			db.ChangeFloor(PeerID, Message)
		}
		return ""
	}
	if Message == "Личный кабинет" {
		createAndSendMessagesAndKeyboard(vk, PeerID, "Чем я могу тебе помочь?", createPersonalAreaKeyboard())
		return Message
	}
	if Message == "Изменить кабинет" {
		db.ChangeRoom(PeerID, "")
		createAndSendMessages(vk, PeerID, "Укажи номер своего кабинета")
		return "Кабинет"

	}
	if Message == "История заказов" {
		sendHistory(vk, port, PeerID)
		createAndSendMessagesAndKeyboard(vk, PeerID, "Выбери с чем тебе помочь", createGeneralKeyboard(true))
		return Message
	}

	if userStatus.LastMessages == "Отменить заказ" {
		userHistory := GetHistory(port, PeerID)

		for _, task := range userHistory.Tasks {
			if strconv.Itoa(int(task.ID)) == Message {
				req, err := http.NewRequest(http.MethodPut, "http://"+port+"/api/task/"+Message+"/4", nil)
				if err != nil {
					fmt.Println(err)
				}
				_, err = http.DefaultClient.Do(req)
				createAndSendMessagesAndKeyboard(vk, PeerID, "Твой заказ отменен", createGeneralKeyboard(false))
				return Message
			}
		}
		createAndSendMessagesAndKeyboard(vk, PeerID, "Этот заказ не может быть отменен", createGeneralKeyboard(false))

		return Message
	}
	if Message == "Вернуться назад" {
		createAndSendMessagesAndKeyboard(vk, PeerID, "Выбери с чем тебе помочь", createGeneralKeyboard(true))
		return Message
	}
	if Message == "Отменить заказ" {
		SelectDeleteHistory(vk, port, PeerID)
		return Message
	}

	if userStatus.LastMessages == "Заказ" && Message != "Сделать заказ" {
		PostNewTask(vk, Message, PeerID, userStatus.Room, userStatus.Floor)
		createAndSendMessagesAndKeyboard(vk, PeerID, "Твой заказ создан: "+Message, createGeneralKeyboard(false))
		return "Заказ создан"

	}

	if Message == "Сделать заказ" {
		createAndSendMessages(vk, PeerID, "Напиши что тебе принести")
		return "Заказ"
	}
	return ""

}
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

func Start(key string, groupId int) {
	vk := api.NewVK(key)
	lp, err := longpoll.NewLongPoll(vk, groupId)
	if err != nil {
		panic(err)
	}

	// Обработка новых сообщений
	lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		Message := obj.Message.Text
		PeerID := obj.Message.PeerID

		log.Printf("New messages: %d:%s", PeerID, Message)

		_, err := db.GetUsers(PeerID)
		if err != nil {
			db.CreateUsers(db.Users{UserID: PeerID, Room: "0"})
			createAndSendMessages(vk, PeerID, "Привет! Я Печенька")
		}

		userFile := messageHandling(vk, Message, PeerID)
		db.ChangeMessage(PeerID, userFile)

	})

	// Запуск
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}

	// Безопасное завершение
	// Ждет пока соединение закроется и события обработаются
	lp.Shutdown()
}
