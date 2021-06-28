package web

import (
	"fmt"

	"ssummers02/Cookies/db"
)

func fillTasksAmounts(page *Page) {
	var err error
	page.OpenTasksAmount, err = db.GetNumberOfOpenTasks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	page.HoldTasksAmount, err = db.GetNumberOfHoldTasks()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
