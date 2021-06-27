package web

import (
	"fmt"
	"net/http"
	"text/template"
)

type Page struct {
	HighlightActiveMenu   string
	HighlightAllMenu      string
	HighlightSettingsMenu string
	OpenTasksAmount       int64
	HoldTasksAmount       int64
	Content               string
}

func AllTasksPage(w http.ResponseWriter, r *http.Request) {
	var page Page
	var err error
	fillTasksAmounts(&page)
	page.HighlightAllMenu = "active"
	page.Content, err = buildTasksTable("all")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	t := template.New("tasks.html")
	t, err = t.ParseFiles("./web/assets/tasks.html")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = t.Execute(w, page)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}

func ActiveTasksPage(w http.ResponseWriter, r *http.Request) {
	var page Page
	var err error
	fillTasksAmounts(&page)
	page.HighlightActiveMenu = "active"
	page.Content, err = buildTasksTable("active")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	t := template.New("tasks.html")
	t, err = t.ParseFiles("./web/assets/tasks.html")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = t.Execute(w, page)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
