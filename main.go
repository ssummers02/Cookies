package main

import (
	"os"
	"ssummers02/Cookies/bot"
	"ssummers02/Cookies/db"
)

func main() {
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0777)
	}
	db.InitDB()
	bot.Start()
}
