package main

import (
<<<<<<< HEAD
	"net/http"
	"ssummers02/Cookies/api"
=======
	"os"
	"ssummers02/Cookies/bot"
>>>>>>> c1cb319ca78afbea37176e031e64350dfbcdf446
	"ssummers02/Cookies/db"
)

func main() {
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0777)
	}
	db.InitDB()
<<<<<<< HEAD
	http.HandleFunc("/", api.GetTasksTable)
	http.HandleFunc("/add_task", api.NewTask)
	http.ListenAndServe(":8080", nil)
=======
	bot.Start()
>>>>>>> c1cb319ca78afbea37176e031e64350dfbcdf446
}
