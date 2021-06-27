package web

import (
	"log"
	"ssummers02/Cookies/db"
	"strconv"
	"strings"
)

func buildTasksTable(room string) (string, error) {
	var tasks []db.Task
	var err error
	switch room {
	case "active":
		tasks, err = db.GetActiveTasks()
		if err != nil {
			log.Printf("%s:%s", "db.GetActiveTasks", err.Error())
			return "", err
		}
	case "all":
		tasks, err = db.GetAllTasks()
		if err != nil {
			log.Printf("%s:%s", "db.GetAllTasks", err.Error())
			return "", err
		}
	default:
		tasks, err = db.GetTaskInRoom(room)
		if err != nil {
			return "", err
		}
	}
	var b strings.Builder
	b.WriteString("<section class=\"content\">\n<div class=\"container-fluid\">\n<div class=\"row\">\n" +
		"<div class=\"col-12\">\n<div class=\"card\">\n<div class=\"card-header\">\n" +
		"<h3 class=\"card-title\">Заявки</h3>\n</div>\n<div class=\"card-body\">\n" +
		"<table id=\"example1\" class=\"table table-bordered table-striped\">\n<thead>\n<tr>\n<th>Заявка</th>\n" +
		"<th>Статус</th>\n<th>Дата</th>\n<th>Этаж</th>\n<th>Кабинет</th>\n<th>Заявка</th>\n<th>Заказчик</th>\n" +
		"<th></th>\n</tr>\n</thead>\n<tbody>")
	for _, v := range tasks {

		statusHTML := ""
		openButtonHTML := "<a class=\"btn btn-info btn-sm\" href=\"#\"><i class=\"fas fa-clipboard-list\"></i>Открыть</a>\n"
		doneButtonHTML := "<a class=\"btn btn-success btn-sm\" href=\"#\"><i class=\"fas fa-check\"></i>Завершить</a>\n"
		holdButtonHTML := "<a class=\"btn btn-warning btn-sm\" href=\"#\"><i class=\"fas fa-question\"></i>Уточнить</a>\n"
		rejectButtonHTML := "<a class=\"btn btn-danger btn-sm\" href=\"#\"><i class=\"fas fa-ban\"></i>Отклонить</a>\n</div></td>\n"
		switch v.Status {
		case 0:
			statusHTML = "info\">Открыта"
			openButtonHTML = "<a class=\"btn btn-outline-light btn-sm\" href=\"#\"><i class=\"fas fa-clipboard-list\"></i>Открыть</a>\n"
		case 1:
			statusHTML = "success\">Выполнена"
			doneButtonHTML = "<a class=\"btn btn-outline-light btn-sm\" href=\"#\"><i class=\"fas fa-check\"></i>Завершить</a>\n"
		case 2:
			statusHTML = "warning\">Требует уточнения"
			holdButtonHTML = "<a class=\"btn btn-outline-light btn-sm\" href=\"#\"><i class=\"fas fa-question\"></i>Уточнить</a>\n"
		case 3:
			statusHTML = "danger\">Отклонена"
			rejectButtonHTML = "<a class=\"btn btn-outline-light btn-sm\" href=\"#\"><i class=\"fas fa-ban\"></i>Отклонить</a>\n</div></td>\n"
		case 4:
			statusHTML = "danger\">Отменена"
			openButtonHTML = "<a class=\"btn btn-outline-light btn-sm\" href=\"#\"><i class=\"fas fa-clipboard-list\"></i>Открыть</a>\n"
			doneButtonHTML = "<a class=\"btn btn-outline-light btn-sm\" href=\"#\"><i class=\"fas fa-check\"></i>Завершить</a>\n"
			holdButtonHTML = "<a class=\"btn btn-outline-light btn-sm\" href=\"#\"><i class=\"fas fa-question\"></i>Уточнить</a>\n"
			rejectButtonHTML = "<a class=\"btn btn-outline-light btn-sm\" href=\"#\"><i class=\"fas fa-ban\"></i>Отклонить</a>\n</div></td>\n"
		}
		b.WriteString("<tr>\n<td>" + strconv.Itoa(int(v.ID)) + "</td>\n" +
			"<td><span class=\"badge badge-" + statusHTML + "</span></td>\n" +
			"<td>" + v.CreatedAt.Format("2006-01-02 15:04") + "</td>\n" +
			"<td>" + strconv.Itoa(v.Floor) + "</td>\n" +
			"<td>" + v.Room + "</td>\n" +
			"<td>" + v.Text + "</td>\n" +
			"<td><a href=\"https://vk.com/id" + strconv.Itoa(int(v.UserID)) +
			"\" target=\"_blank\" rel=\"noopener noreferrer\">" + v.Name + "</a></td>\n" +
			"<td><div class=\"sparkbar\" data-color=\"#00a65a\" data-height=\"20\">\n" +
			openButtonHTML + doneButtonHTML + holdButtonHTML + rejectButtonHTML +
			"</tr>\n")

	}
	b.WriteString("</tbody>\n<tfoot>\n<tr>\n<th>Заявка</th>\n<th>Статус</th>\n<th>Дата</th>\n<th>Этаж</th>\n" +
		"<th>Кабинет</th>\n<th>Заявка</th>\n<th>Заказчик</th>\n<th></th>\n</tr>\n</tfoot>\n</table>\n</div>\n" +
		"</div>\n</div>\n</div>\n</div>\n</section>")
	return b.String(), err
}
