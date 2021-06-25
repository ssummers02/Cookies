package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"ssummers02/Cookies/db"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func messageHandling(vk *api.VK, Message string, PeerID int) string {
	userStatus, _ := db.GetUsers(PeerID)

	port := os.Getenv("ADDRESS")

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
		floor, err := strconv.Atoi(Message)
		if err != nil { // если возникла ошибка
			PostAndSendMessages(vk, PeerID, "Неверный этаж, попробуй еще раз")
		} else {
			PostMessagesAndKeyboard(vk, PeerID, "Твой этаж:"+strconv.Itoa(floor)+"\nЧем я могу тебе помочь?", GetGeneralKeyboard(true))
			db.ChangeFloor(PeerID, Message)
		}
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
		PostHistoryForUser(vk, port, PeerID)
		PostMessagesAndKeyboard(vk, PeerID, "Выбери с чем тебе помочь", GetGeneralKeyboard(true))
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
				PostMessagesAndKeyboard(vk, PeerID, "Твой заказ отменен", GetGeneralKeyboard(false))
				return Message
			}
		}
		PostMessagesAndKeyboard(vk, PeerID, "Этот заказ не может быть отменен", GetGeneralKeyboard(false))

		return Message
	}
	if Message == "Вернуться назад" {
		PostMessagesAndKeyboard(vk, PeerID, "Выбери с чем тебе помочь", GetGeneralKeyboard(true))
		return Message
	}
	if Message == "Отменить заказ" {
		SelectDeleteHistory(vk, port, PeerID)
		return Message
	}

	if userStatus.LastMessages == "Заказ" && Message != "Сделать заказ" {
		PostNewTask(vk, Message, PeerID, userStatus.Room, userStatus.Floor)
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
