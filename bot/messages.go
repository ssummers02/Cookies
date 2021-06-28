package bot

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// Отправка сообщения пользователю
func postAndSendMessages(vk *api.VK, peerId int, text string) {
	rand.Seed(time.Now().UnixNano())
	b := params.NewMessagesSendBuilder()

	b.Message(text)
	b.RandomID(rand.Intn(2147483647))
	b.PeerID(peerId)
	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "messages",
			"func":    "postAndSendMessages",
			"error":   err,
		}).Warning("error send messages")
	}
}

// Отправка сообщения пользователю
func postMessagesAndKeyboard(vk *api.VK, peerId int, text string, k *object.MessagesKeyboard) {
	rand.Seed(time.Now().UnixNano())
	b := params.NewMessagesSendBuilder()
	b.Keyboard(k)
	b.Message(text)
	b.RandomID(rand.Intn(2147483647))
	b.PeerID(peerId)
	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "messages",
			"func":    "postMessagesAndKeyboard",
			"error":   err,
		}).Warning("error send messages")
	}
}
func postMessageAdm(vk *api.VK, message string, room string, floor int) {
	adm, _ := strconv.Atoi(os.Getenv("ADM"))
	res := "Новый заказ\n" +
		"Этаж: " + strconv.Itoa(floor) +
		"\nКабинет: " + room +
		"\n" + message
	postAndSendMessages(vk, adm, res)
}
