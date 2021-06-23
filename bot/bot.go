package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
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
		user := OpenUserFile(nameFile)
		log.Print(user)

		if user.Cabinet == 0 && user.LastMessages == "Кабинет" {
			cab, err := strconv.Atoi(Message)
			if err != nil { // если возникла ошибка
				createAndSendMessages(vk, PeerID, "Неверный кабинет, попробуй еще раз")
			} else {
				createAndSendMessagesAndKeyboard(vk, PeerID, "Твой новый кабинет:"+Message, createGeneralKeyboard(true))
				user.Cabinet = cab
			}
		}

		if user.Cabinet == 0 && user.LastMessages != "Кабинет" {
			createAndSendMessages(vk, PeerID, "Я тебя не знаю, давай познакомимься поближе\nУкажи номер своего кабинета")
			user.LastMessages = "Кабинет"
		}
		if Message == "Личный кабинет" {
			user.LastMessages = Message
			createAndSendMessagesAndKeyboard(vk, PeerID, "Чем я могу тебе помочь?", createPersonalAreaKeyboard())
		}
		if Message == "Изменить кабинет" {
			user.Cabinet = 0
			user.LastMessages = "Кабинет"
			createAndSendMessages(vk, PeerID, "Укажи номер своего кабинета")
		}
		if Message == "История заказов" {
			user.LastMessages = Message
			// отправка истории 5 штук
			createAndSendMessagesAndKeyboard(vk, PeerID, "ТУТ ИСТОРИЯ", createGeneralKeyboard(false))

		}

		if user.LastMessages == "Отменить заказ" {
			user.LastMessages = Message
			// получаем id и отменяем заказ
			createAndSendMessagesAndKeyboard(vk, PeerID, "Твой заказ отменен", createGeneralKeyboard(false))

		}

		if Message == "Отменить заказ" {
			user.LastMessages = Message
			// отправка истории 5 штук
			/*			createAndSendMessagesAndKeyboard(vk, PeerID, "Выбери заказ", createPersonalAreaKeyboard())
			 */createAndSendMessages(vk, PeerID, "ТУТ ИСТОРИЯ с inline кнопками")

		}

		if user.LastMessages == "Заказ" {
			user.LastMessages = "Заказ создан"
			// создать заявку
			createAndSendMessagesAndKeyboard(vk, PeerID, "Твой заказ создан: "+Message, createGeneralKeyboard(false))

		}

		if Message == "Сделать заказ" {
			user.LastMessages = "Заказ"
			createAndSendMessages(vk, PeerID, "Напиши что тебе принести")
		}

		changeUserFile(nameFile, user)

	})

	// Запуск
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}

	// Безопасное завершение
	// Ждет пока соединение закроется и события обработаются
	lp.Shutdown()
}
