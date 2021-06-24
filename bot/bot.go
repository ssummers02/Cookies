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
	"ssummers02/Cookies/db"
	"strconv"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/SevereCloud/vksdk/v2/object"
)

type Users struct {
	LastMessages string `json:"LastMessages"`
	Cabinet      int    `json:"Cabinet"`
}

type ArrayTask struct {
	Tasks []db.Task
}

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

// Создание файла пользователя
func createNewUser(nameFile string) {
	file, err := os.Create(nameFile) // создаем файл

	if err != nil { // если возникла ошибка
		log.Print("Unable to create file:", err)
	}
	emp := &Users{"", 0} // default значения
	e, err := json.Marshal(emp)
	file.WriteString(string(e))

	defer file.Close()
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
func OpenUserFile(nameFile string) Users {
	var user Users

	jsonFile, err := os.Open(nameFile) // Открытие jsonFile
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile) //Считывание и раскодирование в json
	json.Unmarshal(byteValue, &user)
	return user

}
func changeUserFile(nameFile string, users Users) {
	file, err := os.Create(nameFile) // создаем файл

	if err != nil { // если возникла ошибка
		log.Print("Unable to create file:", err)
	}

	e, err := json.Marshal(users)
	file.WriteString(string(e))

	defer file.Close()
}

func messageHandling(vk *api.VK, userStatus Users, Message string, PeerID int) Users {

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
		// отправка истории 5 штук

		resp, err := http.Get("http://" + port + "/user/" + strconv.Itoa(PeerID) + "/5")

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}
		var userHistory ArrayTask
		json.Unmarshal(body, &userHistory)
		log.Print(userHistory)
		createMessage := ""
		for i := 0; i < len(userHistory.Tasks); i++ {
			createMessage = createMessage + "№" + strconv.Itoa(int(userHistory.Tasks[i].ID)) + ": " + userHistory.Tasks[i].Text + "-" + findOutTheStatus(userHistory.Tasks[i].Status) + "\n"
		}
		createAndSendMessagesAndKeyboard(vk, PeerID, createMessage, createGeneralKeyboard(false))
		return userStatus
	}

	if userStatus.LastMessages == "Отменить заказ" {
		userStatus.LastMessages = Message
		// получаем id и отменяем заказ
		createAndSendMessagesAndKeyboard(vk, PeerID, "Твой заказ отменен", createGeneralKeyboard(false))
		return userStatus
	}

	if Message == "Отменить заказ" {
		userStatus.LastMessages = Message
		// отправка истории 5 штук
		/*			createAndSendMessagesAndKeyboard(vk, PeerID, "Выбери заказ", createPersonalAreaKeyboard())
		 */createAndSendMessages(vk, PeerID, "ТУТ ИСТОРИЯ с inline кнопками")
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

	//Обработка новых сообщений
	lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		Message := obj.Message.Text
		PeerID := obj.Message.PeerID

		log.Printf("New messages: %d:%s", PeerID, Message)

		nameFile := "temp/" + strconv.Itoa(PeerID) + ".json"
		if !Exists(nameFile) {
			createNewUser(nameFile)
			createAndSendMessages(vk, PeerID, "Привет! Я Печенька")

		}
		userFile := OpenUserFile(nameFile)
		userFile = messageHandling(vk, userFile, Message, PeerID)
		changeUserFile(nameFile, userFile)

	})

	// Запуск
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}

	// Безопасное завершение
	// Ждет пока соединение закроется и события обработаются
	lp.Shutdown()
}
