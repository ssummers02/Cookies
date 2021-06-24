package bot

import (
	"encoding/json"
	"log"
	"os"

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

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
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
