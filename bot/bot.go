package bot

import (
	"context"
	"log"
	"os"

	"ssummers02/Cookies/db"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

var port string

func messageHandling(vk *api.VK, Message string, PeerID int) string {
	userStatus, _ := db.GetUsers(PeerID)

	if userStatus.Room == "0" && userStatus.LastMessages == "Кабинет" {
		db.ChangeRoom(PeerID, Message)
		PostAndSendMessages(vk, PeerID, "Твой новый кабинет: "+Message+"\n Укажи этаж")
		return "Этаж"
	}
	if userStatus.Room == "0" && userStatus.LastMessages != "Кабинет" {
		PostAndSendMessages(vk, PeerID, "Я тебя не знаю, давай познакомимься поближе\nУкажи номер своего кабинета")
		return "Кабинет"
	}
	if userStatus.LastMessages == "Этаж" {
		PostFloor(vk, Message, PeerID)
		return ""
	}
	if Message == "Личный кабинет" {
		PostMessagesAndKeyboard(vk, PeerID, "Чем я могу тебе помочь?", GetPersonalAreaKeyboard())
		return Message
	}
	if Message == "Изменить кабинет" {
		db.ChangeRoom(PeerID, "")
		PostAndSendMessages(vk, PeerID, "Укажи номер своего кабинета")
		return "Кабинет"
	}
	if Message == "История заказов" {
		PostHistoryForUser(vk, PeerID)
		PostMessagesAndKeyboard(vk, PeerID, "Выбери с чем тебе помочь", GetGeneralKeyboard(true))
		return Message
	}
	if userStatus.LastMessages == "Отменить заказ" {
		return ChangeStatus(vk, Message, PeerID)
	}
	if Message == "Вернуться назад" {
		PostMessagesAndKeyboard(vk, PeerID, "Выбери с чем тебе помочь", GetGeneralKeyboard(true))
		return Message
	}
	if Message == "Отменить заказ" {
		SelectDeleteHistory(vk, PeerID)
		return Message
	}
	if userStatus.LastMessages == "Заказ" && Message != "Сделать заказ" {
		PostNewTask(vk, Message, PeerID, userStatus.Room, userStatus.Floor)
		postMessageAdm(vk, Message, userStatus.Room, userStatus.Floor)
		PostMessagesAndKeyboard(vk, PeerID, "Твой заказ создан: "+Message, GetGeneralKeyboard(false))
		return "Заказ создан"
	}
	if Message == "Сделать заказ" {
		PostAndSendMessages(vk, PeerID, "Напиши что тебе принести")
		return "Заказ"
	}
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
		Message := obj.Message.Text
		PeerID := obj.Message.PeerID

		log.Printf("New messages: %d:%s", PeerID, Message)

		_, err := db.GetUsers(PeerID)
		if err != nil {
			db.CreateUsers(db.Users{UserID: PeerID, Room: "0"})
			PostAndSendMessages(vk, PeerID, "Привет! Я Печенька")
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
