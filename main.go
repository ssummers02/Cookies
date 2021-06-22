package main

import (
	"net/http"
	"ssummers02/Cookies/api"
	"ssummers02/Cookies/db"
)

func main() {
	db.InitDB()
	http.HandleFunc("/", api.GetTasksTable)
	http.HandleFunc("/add_task", api.NewTask)
	http.ListenAndServe(":8080", nil)
}
