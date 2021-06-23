package bot

import (
	"github.com/SevereCloud/vksdk/v2/object"
	"os"
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

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func createDelete(ar []int) *object.MessagesKeyboard {
	k := object.NewMessagesKeyboardInline()
	for _, value := range ar {
		k.AddRow()
		k.AddTextButton(string(rune(value)), ``, `primary`)
	}

	return k
}
