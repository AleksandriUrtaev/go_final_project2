package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/AleksandriUrtaev/go_final_project2/pkg/db"
	"github.com/AleksandriUrtaev/go_final_project2/pkg/scheduler"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		} else {
			log.Println("taskDone GetTask error:", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	// одноразовая задача → удаляем
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			log.Println("taskDone DeleteTask error:", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Не удалось удалить задачу"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]interface{}{})
		return
	}

	// периодическая задача → вычисляем следующую дату
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	next, err := scheduler.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Некорректное правило повторения"})
		return
	}

	if err := db.UpdateDate(id, next); err != nil {
		log.Println("taskDone UpdateDate error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Не удалось обновить дату задачи"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
