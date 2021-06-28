package web

import (
	"log"
	"ssummers02/Cookies/db"
)

func fillTasksAmounts(page *Page) {
	var err error
	page.OpenTasksAmount, err = db.GetNumberOfOpenTasks()
	if err != nil {
		log.Print(err)
		return
	}
	page.HoldTasksAmount, err = db.GetNumberOfHoldTasks()
	if err != nil {
		log.Print(err)
		return
	}
}
