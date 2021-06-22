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
	router := http.NewServeMux()
	router.HandleFunc("/", api.GetTasksTable)
	router.HandleFunc("/add_task", api.NewTask)
	router.HandleFunc("/users/{id}", users.showHandler)
	http.ListenAndServe(":8080", nil)
	bot.Start()
}
