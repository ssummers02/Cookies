package bot

import "github.com/SevereCloud/vksdk/v2/object"

func getPersonalAreaKeyboard() *object.MessagesKeyboard {
	k := object.NewMessagesKeyboardInline()

	k.AddRow()
	k.AddTextButton(`Изменить кабинет`, ``, `primary`)

	k.AddRow()
	k.AddTextButton(`История заказов`, ``, `secondary`)

	k.AddRow()
	k.AddTextButton(`Отменить заказ`, ``, `negative`)

	return k
}
func getGeneralKeyboard(t bool) *object.MessagesKeyboard {
	k := object.NewMessagesKeyboard(object.BaseBoolInt(t))

	k.AddRow()
	k.AddTextButton(`Личный кабинет`, ``, `primary`)
	k.AddTextButton(`Сделать заказ`, ``, `primary`)

	return k
}
