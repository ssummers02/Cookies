package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"ssummers02/Cookies/api"
	"ssummers02/Cookies/bot"
	"ssummers02/Cookies/db"
	"ssummers02/Cookies/web"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("ADDRESS")
	dbName := os.Getenv("DATABASE")
	limit := os.Getenv("LIMIT")
	vkKey := os.Getenv("VK_KEY")
	vkGroup := os.Getenv("VK_GROUP_ID")
	/*	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		os.Mkdir("temp", 0777)
	}*/
	lim, _ := strconv.Atoi(limit)
	vkGroupId, _ := strconv.Atoi(vkGroup)
	go func(vkKey string, vkGroupId int) {
		bot.Start(vkKey, vkGroupId)
	}(vkKey, vkGroupId)
	db.InitDB(dbName, lim)
	router := mux.NewRouter()
	router.HandleFunc("/api/", api.GetTasksTable).Methods("GET")
	router.HandleFunc("/api/add_task", api.NewTask).Methods("POST")
	router.HandleFunc("/api/room/{id:[0-9]+}", api.GetTasksInRoom).Methods("GET")
	router.HandleFunc("/api/task/{id:[0-9]+}", api.DeleteTask).Methods("DELETE")
	router.HandleFunc("/api/user/{user_id:[0-9]+}/{count:[0-9]+}", api.GetHistory).Methods("GET")
	router.HandleFunc("/api/useractive/{user_id:[0-9]+}/{count:[0-9]+}", api.GetActiveHistory).Methods("GET")
	router.HandleFunc("/api/task/status/{task_id:[0-9]+}/{status:[0-9]+}", api.ChangeStatus).Methods("PUT")
	router.HandleFunc("/api/task/floor/{user_id:[0-9]+}/{floor:[0-9]+}", api.ChangeFloor).Methods("PUT")

	// router.HandleFunc("/", web.ShowActiveTasks).Methods("GET")
	// router.HandleFunc("/all", web.ShowAllTasks).Methods("GET")
	// router.HandleFunc("/settings", web.ShowSettings).Methods("GET")

	router.HandleFunc("/", web.ActiveTasksPage).Methods("GET")
	router.HandleFunc("/alltasks", web.AllTasksPage).Methods("GET")
	router.PathPrefix("/plugins/").Handler(http.StripPrefix("/plugins/", http.FileServer(http.Dir("./web/assets/plugins"))))
	router.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("./web/assets/dist"))))

	srv := &http.Server{
		Handler: router,
		Addr:    port,
	}
	log.Fatal(srv.ListenAndServe())
}
