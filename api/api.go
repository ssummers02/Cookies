package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"ssummers02/Cookies/db"
	"strconv"

	"github.com/gorilla/mux"
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

func GetTasksInRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room_id := vars["id"]
	tasks, err := db.GetTaskInRoom(room_id)
	if err != nil {
		fmt.Println(err.Error())
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
	var task db.Task
	byteVale, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.Unmarshal(byteVale, &task)
	if err := db.CreateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	task_id := vars["id"]
	taskId, _ := strconv.Atoi(task_id)
	if err := db.DeleteTask(uint(taskId)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// PutTask: PUT handler for task model.
func PutTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task db.Task
	byteVale, err := ioutil.ReadAll(r.Body)

	json.Unmarshal(byteVale, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	countVar := vars["count"]
	userId, _ := strconv.Atoi(userID)
	count, _ := strconv.Atoi(countVar)
	tasks, err := db.GetUserHistory(uint(userId), count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newStatus := vars["status"]
	taskId := vars["task_id"]
	if err := db.ChangeStatus(taskId, newStatus); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
