package bot

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

// Отправка сообщения пользователю
func PostAndSendMessages(vk *api.VK, PeerID int, text string) {
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
func PostMessagesAndKeyboard(vk *api.VK, PeerID int, text string, k *object.MessagesKeyboard) {
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
func postMessageAdm(vk *api.VK, Message string, room string, floor int) {
	adm, _ := strconv.Atoi(os.Getenv("ADM"))
	res := "Новый заказ\n" +
		"Этаж: " + strconv.Itoa(floor) +
		"\nКабинет" + room +
		"\n" + Message
	PostAndSendMessages(vk, adm, res)
}
