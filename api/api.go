package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"ssummers02/Cookies/bot"
	"ssummers02/Cookies/db"

	"github.com/gorilla/mux"
)

type Response struct {
	Tasks []db.Task
}

func GetTasksTable(w http.ResponseWriter) {
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
	_, err = w.Write(js)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "api",
			"func":    "GetTasksTable",
			"error":   err,
		}).Warning("err Get Tasks Table")
	}
}

func GetTasksInRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["id"]
	tasks, err := db.GetTaskInRoom(roomId)
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
	_, err = w.Write(js)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "api",
			"func":    "GetTasksInRoom",
			"error":   err,
		}).Warning("err Get Tasks In Room")
	}
}

func NewTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task db.Task
	byteVale, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(byteVale, &task)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "api",
			"func":    "NewTask",
			"error":   err,
		}).Warning("err create New Task")
	}
	if err := db.CreateTask(task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, _ := strconv.Atoi(vars["id"])
	if err := db.DeleteTask(uint(taskId)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func GetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, _ := strconv.Atoi(vars["user_id"])
	count, _ := strconv.Atoi(vars["count"])
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
	_, err = w.Write(js)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "api",
			"func":    "GetHistory",
			"error":   err,
		}).Warning("err Get History")
	}
}

func GetActiveHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	countVar := vars["count"]
	userId, _ := strconv.Atoi(userID)
	count, _ := strconv.Atoi(countVar)
	tasks, err := db.GetUserActiveHistory(uint(userId), count)
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
	_, err = w.Write(js)
	if err != nil {
		log.WithFields(log.Fields{
			"package": "api",
			"func":    "GetActiveHistory",
			"error":   err,
		}).Warning("err Get Active History")
	}
}

func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newStatus := vars["status"]
	taskId := vars["task_id"]
	if err := db.ChangeStatus(taskId, newStatus); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	taskID, _ := strconv.Atoi(taskId)
	task, _ := db.GetTask(uint(taskID)) // Ошибки не может быть ибо существование task с таким taskID проверенно выше
	st, _ := strconv.Atoi(newStatus)
	temp := db.StatusChangeAlert{RecipientUserID: int(task.UserID), TaskID: taskId, TaskText: task.Text, Status: bot.FindOutTheStatus(uint(st))}
	bot.PostChangeStatus(temp)
}

func ChangeFloor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newFloor, _ := strconv.Atoi(vars["floor"])
	userId, _ := strconv.Atoi(vars["user_id"])
	if err := db.ChangeFloor(userId, newFloor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
