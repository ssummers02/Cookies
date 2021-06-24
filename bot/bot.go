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
	}
	return ""
}

func OpenUserFile(nameFile string) db.Users {
	var user db.Users

	jsonFile, err := os.Open(nameFile) // Открытие jsonFile
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile) // Считывание и раскодирование в json
	json.Unmarshal(byteValue, &user)
	return user

}
func GetHistory(port string, PeerID int) db.ArrayTask {
	var userHistory db.ArrayTask

	resp, err := http.Get("http://" + port + "/user/" + strconv.Itoa(PeerID) + "/5")

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

func messageHandling(vk *api.VK, Message string, PeerID int) db.Users {
	userStatus, _ := db.GetUsers(PeerID)

	port := os.Getenv("ADDRESS")

	if userStatus.Cabinet == 0 && userStatus.LastMessages == "Кабинет" {
		cab, err := strconv.Atoi(Message)
		if err != nil { // если возникла ошибка
			createAndSendMessages(vk, PeerID, "Неверный кабинет, попробуй еще раз")
		} else {
			createAndSendMessagesAndKeyboard(vk, PeerID, "Твой новый кабинет:"+Message, createGeneralKeyboard(true))
			userStatus.Cabinet = cab
		}
		return userStatus
	}

	if userStatus.Cabinet == 0 && userStatus.LastMessages != "Кабинет" {
		createAndSendMessages(vk, PeerID, "Я тебя не знаю, давай познакомимься поближе\nУкажи номер своего кабинета")
		userStatus.LastMessages = "Кабинет"
		return userStatus

	}
	if Message == "Личный кабинет" {
		userStatus.LastMessages = Message
		createAndSendMessagesAndKeyboard(vk, PeerID, "Чем я могу тебе помочь?", createPersonalAreaKeyboard())
		return userStatus
	}
	if Message == "Изменить кабинет" {
		userStatus.Cabinet = 0
		userStatus.LastMessages = "Кабинет"
		createAndSendMessages(vk, PeerID, "Укажи номер своего кабинета")
		return userStatus

	}
	if Message == "История заказов" {
		userStatus.LastMessages = Message
		sendHistory(vk, port, PeerID)
		createAndSendMessagesAndKeyboard(vk, PeerID, "Выбери с чем тебе помочь", createGeneralKeyboard(true))
		return userStatus
	}

	if userStatus.LastMessages == "Отменить заказ" {
		userStatus.LastMessages = Message

		req, err := http.NewRequest(http.MethodDelete, "http://"+port+"/task/"+Message, nil)
		if err != nil {
			fmt.Println(err)
		}
		_, err = http.DefaultClient.Do(req)

		createAndSendMessagesAndKeyboard(vk, PeerID, "Твой заказ отменен", createGeneralKeyboard(false))
		return userStatus
	}
	if Message == "Вернуться назад" {
		userStatus.LastMessages = Message
		createAndSendMessagesAndKeyboard(vk, PeerID, "Выбери с чем тебе помочь", createGeneralKeyboard(true))
		return userStatus
	}
	if Message == "Отменить заказ" {
		userStatus.LastMessages = Message
		SelectDeleteHistory(vk, port, PeerID)
		return userStatus
	}

	if userStatus.LastMessages == "Заказ" && Message != "Сделать заказ" {
		userStatus.LastMessages = "Заказ создан"
		// создать заявку
		emp := &db.Task{UserID: uint(PeerID), Room: uint(userStatus.Cabinet), Text: Message} // default значения
		jsonData, _ := json.Marshal(emp)

		_, err := http.Post("http://"+port+"/add_task", "application/json",
			bytes.NewBuffer(jsonData))

		if err != nil {
			log.Fatal(err)
		}

		createAndSendMessagesAndKeyboard(vk, PeerID, "Твой заказ создан: "+Message, createGeneralKeyboard(false))
		return userStatus

	}

	if Message == "Сделать заказ" {
		userStatus.LastMessages = "Заказ"
		createAndSendMessages(vk, PeerID, "Напиши что тебе принести")
		return userStatus
	}
	return userStatus

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
			db.CreateUsers(db.Users{UserID: PeerID})
			createAndSendMessages(vk, PeerID, "Привет! Я Печенька")
		}

		userFile := messageHandling(vk, Message, PeerID)
		db.UpdateUsers(userFile)

	})

	// Запуск
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}

	// Безопасное завершение
	// Ждет пока соединение закроется и события обработаются
	lp.Shutdown()
}
