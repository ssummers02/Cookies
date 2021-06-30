package web

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"ssummers02/Cookies/db"
	"text/template"
)

type Page struct {
	HighlightActiveMenu   bool
	HighlightAllMenu      bool
	HighlightSettingsMenu bool
	OpenTasksAmount       int64
	HoldTasksAmount       int64
	Content               string
}

func AllTasksPage(w http.ResponseWriter, r *http.Request) {
	var page Page
	var err error
	fillTasksAmounts(&page)
	page.HighlightAllMenu = true
	page.Content, err = buildTasksTable("all")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	t := template.New("main.html")
	t, err = t.ParseFiles("./web/assets/templates/main.html")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = t.Execute(w, page)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func ActiveTasksPage(w http.ResponseWriter, r *http.Request) {
	var page Page
	var err error
	fillTasksAmounts(&page)
	page.HighlightActiveMenu = true
	page.Content, err = buildTasksTable("active")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	t := template.New("main.html")
	t, err = t.ParseFiles("./web/assets/templates/main.html")
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = t.Execute(w, page)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func ChangeStatusOnActiveTasksPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := db.ChangeStatus(vars["taskID"], vars["statusID"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ActiveTasksPage(w, r)
}

func ChangeStatusOnAllTasksPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if err := db.ChangeStatus(vars["taskID"], vars["statusID"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	AllTasksPage(w, r)
}
