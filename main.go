package main

import (
	"log"
	"net/http"
	"os"
	"ssummers02/Cookies/api"
	"ssummers02/Cookies/bot"
	"ssummers02/Cookies/db"

	"github.com/gorilla/mux"
)

func main() {
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0777)
	}
	db.InitDB()
	router := mux.NewRouter()
	router.HandleFunc("/", api.GetTasksTable)
	router.HandleFunc("/add_task", api.NewTask)
	router.HandleFunc("/room/{id}", api.GetTasksInRoom)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
	}
	log.Fatal(srv.ListenAndServe())
	bot.Start()
}
