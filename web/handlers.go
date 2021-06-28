package web

import (
	"fmt"
	"net/http"
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
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	t := template.New("main.html")
	t, err = t.ParseFiles("./web/assets/templates/main.html")
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
	page.HighlightActiveMenu = true
	page.Content, err = buildTasksTable("active")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	t := template.New("main.html")
	t, err = t.ParseFiles("./web/assets/templates/main.html")
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
