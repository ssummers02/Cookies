package web

import (
	"bytes"
	"html"
	"log"
	"ssummers02/Cookies/db"
	"text/template"
)

func buildTasksTable(room string) (string, error) {
	var tasks []db.Task
	var err error
	var renderedTemplate bytes.Buffer
	switch room {
	case "active":
		tasks, err = db.GetActiveTasks()
		if err != nil {
			log.Printf("%s:%s", "buildTasksTable.go - db.GetActiveTasks", err.Error())
			return "", err
		}
	case "all":
		tasks, err = db.GetAllTasks()
		if err != nil {
			log.Printf("%s:%s", "buildTasksTable - db.GetAllTasks", err.Error())
			return "", err
		}
	default:
		tasks, err = db.GetTaskInRoom(room)
		if err != nil {
			return "", err
		}
	}
	t := template.New("tasksTable.html").Funcs(template.FuncMap{"escape": escapingInTemplate})
	t, err = t.ParseFiles("./web/assets/templates/tasksTable.html")
	if err != nil {
		log.Printf("%s:%s", "buildTasksTable t.ParseFiles", err.Error())
		return "", err
	}
	err = t.Execute(&renderedTemplate, tasks)
	if err != nil {
		log.Printf("%s:%s", "buildTasksTable t.Execute", err.Error())
		return "", err
	}
	return renderedTemplate.String(), err
}

func escapingInTemplate(str string) string {
	return html.EscapeString(str)
}
