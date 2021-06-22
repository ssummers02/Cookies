package main

import (
	"ssummers02/Cookies/bot"
	"ssummers02/Cookies/db"
)

func main() {
	db.InitDB()
	bot.Start()
}
