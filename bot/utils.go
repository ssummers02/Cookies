package bot

import (
	"log"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/object"
)

func createPersonalAreaKeyboard() *object.MessagesKeyboard {
	k := object.NewMessagesKeyboardInline()

	k.AddRow()
	k.AddTextButton(`Изменить кабинет`, ``, `primary`)

	k.AddRow()
	k.AddTextButton(`История заказов`, ``, `secondary`)

	k.AddRow()
	k.AddTextButton(`Отменить заказ`, ``, `negative`)

	return k
}
func createGeneralKeyboard(t bool) *object.MessagesKeyboard {
	k := object.NewMessagesKeyboard(object.BaseBoolInt(t))

	k.AddRow()
	k.AddTextButton(`Личный кабинет`, ``, `primary`)
	k.AddTextButton(`Сделать заказ`, ``, `primary`)

	return k
}
func GetName(vk *api.VK, PeerID int) string {
	b := params.NewUsersGetBuilder()
	var id = []string{strconv.Itoa(PeerID)}
	b.UserIDs(id)

	resp, err := vk.UsersGet(b.Params)

	if err != nil {
		log.Fatal(err)
	}
	return resp[0].FirstName + " " + resp[0].LastName
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
	case 4:
		return "отменена пользователем"
	}
	return ""
}
