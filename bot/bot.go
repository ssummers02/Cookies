package bot

import (
	"context"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"math/rand"
)

func createAndSendMessages(vk *api.VK, PeerID int, text string) {
	b := params.NewMessagesSendBuilder()
	b.Message(text)
	b.RandomID(int(rand.Int31()))
	b.PeerID(PeerID)
	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.Fatal(err)
	}
}
func Start() {
	vk := api.NewVK("aa6f5be89eb316d1fbdfb1fab2d82a8229aec785fa980bdc51d606e03b36b1ec5f740cee6cd56d132efce")
	lp, err := longpoll.NewLongPoll(vk, 204006771)
	if err != nil {
		panic(err)
	}

	lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		log.Printf("%d:%s", obj.Message.PeerID, obj.Message.Text)
		PeerID := obj.Message.PeerID

		if obj.Message.Text == "ping" {
			createAndSendMessages(vk, PeerID, "tesst")
		}

	})

	// Запуск
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}

	// Безопасное завершение
	// Ждет пока соединение закроется и события обработаются
	lp.Shutdown()
}
