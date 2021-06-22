package main

import (
	"net/http"
	"os"
	"ssummers02/Cookies/api"
	"ssummers02/Cookies/bot"
	"ssummers02/Cookies/db"
)

func main() {
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0777)
	}
	db.InitDB()
	http.HandleFunc("/", api.GetTasksTable)
	http.HandleFunc("/add_task", api.NewTask)
	http.ListenAndServe(":8080", nil)
	bot.Start()
}
