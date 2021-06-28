package web

import (
	"bytes"
	"html"
	"log"
	"text/template"

	"ssummers02/Cookies/db"
)

func buildTasksTable(room string) (string, error) {
	var tasks []db.Task
	var err error
	var renderedTemplate bytes.Buffer
	switch room {
	case "active":
		tasks, err = db.GetActiveTasks()
		if err != nil {
			log.Print(err)
			return "", err
		}
	case "all":
		tasks, err = db.GetAllTasks()
		if err != nil {
			log.Print(err)
			return "", err
		}
	default:
		tasks, err = db.GetTaskInRoom(room)
		if err != nil {
			log.Print(err)
			return "", err
		}
	}
	t := template.New("tasksTable.html").Funcs(template.FuncMap{"escape": escapingInTemplate})
	t, err = t.ParseFiles("./web/assets/templates/tasksTable.html")
	if err != nil {
		log.Print(err)
		return "", err
	}
	err = t.Execute(&renderedTemplate, tasks)
	if err != nil {
		log.Print(err)
		return "", err
	}
	return renderedTemplate.String(), err
}

func escapingInTemplate(str string) string {
	return html.EscapeString(str)
}
