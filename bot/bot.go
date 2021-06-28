package bot

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"ssummers02/Cookies/db"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

var port string

func messageHandling(vk *api.VK, message string, peerId int) string {
	userStatus, _ := db.GetUsers(peerId)

	if userStatus.Room == "0" && userStatus.LastMessages == "Кабинет" {
		db.ChangeRoom(peerId, message)
		postAndSendMessages(vk, peerId, "Твой новый кабинет: "+message+"\nУкажи этаж")
		return "Этаж"
	}
	if userStatus.Room == "0" && userStatus.LastMessages != "Кабинет" {
		postAndSendMessages(vk, peerId, "Я тебя не знаю, давай познакомимься поближе\nУкажи номер своего кабинета")
		return "Кабинет"
	}
	if userStatus.LastMessages == "Этаж" {
		return postFloor(vk, message, peerId)
	}
	if message == "Личный кабинет" {
		postMessagesAndKeyboard(vk, peerId, "Чем я могу тебе помочь?", getPersonalAreaKeyboard())
		return message
	}
	if message == "Изменить кабинет" {
		db.ChangeRoom(peerId, "0")
		db.ChangeFloor(peerId, 0)
		postAndSendMessages(vk, peerId, "Укажи номер своего кабинета")
		return "Кабинет"
	}
	if message == "История заказов" {
		postHistoryForUser(vk, peerId)
		postMessagesAndKeyboard(vk, peerId, "Выбери с чем тебе помочь", getGeneralKeyboard(true))
		return message
	}
	if message == "Вернуться назад" {
		postMessagesAndKeyboard(vk, peerId, "Выбери с чем тебе помочь", getGeneralKeyboard(true))
		return message
	}
	if userStatus.LastMessages == "Отменить заказ" {
		return changeStatus(vk, message, peerId)
	}
	if message == "Отменить заказ" {
		selectDeleteHistory(vk, peerId)
		return message
	}
	if userStatus.LastMessages == "Заказ" && message != "Сделать заказ" {
		postNewTask(vk, message, peerId, userStatus.Room, userStatus.Floor)
		postMessageAdm(vk, message, userStatus.Room, userStatus.Floor)
		postMessagesAndKeyboard(vk, peerId, "Твой заказ создан: "+message, getGeneralKeyboard(false))
		return "Заказ создан"
	}
	if message == "Сделать заказ" {
		postAndSendMessages(vk, peerId, "Напиши что тебе принести")
		return "Заказ"
	}
	postMessagesAndKeyboard(vk, peerId, "Я всего лишь печенька и не знаю такого, попробуй еще раз", getGeneralKeyboard(false))

	return ""
}

func Start(key string, groupId int) {
	vk := api.NewVK(key)
	lp, err := longpoll.NewLongPoll(vk, groupId)
	if err != nil {
		panic(err)
	}
	port = os.Getenv("ADDRESS")

	// Обработка новых сообщений
	lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		message := obj.Message.Text
		peerId := obj.Message.PeerID

		log.Printf("New messages: %d:%s", peerId, message)

		_, err := db.GetUsers(peerId)
		if err != nil {
			db.CreateUsers(db.Users{UserID: peerId, Room: "0"})
			postAndSendMessages(vk, peerId, "Привет! Я Печенька")
		}

		userFile := messageHandling(vk, message, peerId)
		db.ChangeMessage(peerId, userFile)

	})

	// Запуск
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}

	// Безопасное завершение
	// Ждет пока соединение закроется и события обработаются
	lp.Shutdown()
}
