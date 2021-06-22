package api

import (
	"encoding/json"
	"net/http"
	"ssummers02/Cookies/db"
)

type Response struct {
	Tasks []db.Task
}

func GetTasksTable(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	resp := Response{Tasks: tasks}
	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func NewTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var task db.Task
	if err := decoder.Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.CreateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
