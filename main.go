package main

import (
	"log"
	"net/http"
	"os"
	"ssummers02/Cookies/api"
	"ssummers02/Cookies/bot"
	"ssummers02/Cookies/db"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("ADDRESS")
	limit := os.Getenv("LIMIT")
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0777)
	}
	lim, err := strconv.Atoi(limit)
	if err != nil {
		log.Fatal(err)
		return
	}
	db.InitDB(lim)
	router := mux.NewRouter()
	router.HandleFunc("/", api.GetTasksTable).Methods("GET")
	router.HandleFunc("/add_task", api.NewTask).Methods("POST")
	router.HandleFunc("/room/{id}", api.GetTasksInRoom).Methods("GET")
	router.HandleFunc("/task/{id}", api.DeleteTask).Methods("DELETE")
	router.HandleFunc("/task/{id}", api.PutTask).Methods("PUT")

	srv := &http.Server{
		Handler: router,
		Addr:    port,
	}
	log.Fatal(srv.ListenAndServe())
	bot.Start()
}
