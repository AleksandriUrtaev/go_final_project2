package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksandriUrtaev/go_final_project2/pkg/db"
	"github.com/AleksandriUrtaev/go_final_project2/pkg/scheduler"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		updateTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// 1. JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Некорректный JSON",
		})
		return
	}

	// 2. title обязателен
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан заголовок задачи",
		})
		return
	}

	// 3. проверка даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	// 4. добавляем в БД
	id, err := db.AddTask(&task)
	if err != nil {
		log.Println("AddTask error:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Ошибка при добавлении задачи",
		})
		return
	}

	// 5. успех
	writeJSON(w, http.StatusCreated, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("writeJSON error:", err)
	}
}

func checkDate(task *db.Task) error {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	var next string

	if task.Repeat != "" {
		next, err = scheduler.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format(dateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан идентификатор",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{
				"error": "Задача не найдена",
			})
		} else {
			log.Println("getTask error:", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "Внутренняя ошибка сервера",
			})
		}
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// 1. Декодирование JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// 2. Проверка ID
	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// 3. Проверка заголовка
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// 4. Проверка даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// 5. Обновление в БД
	if err := db.UpdateTask(&task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		} else {
			log.Println("updateTask error:", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	// 6. Успешный ответ
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	// 1. Проверка ID
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// 2. Удаление задачи
	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		} else {
			log.Println("deleteTask error:", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Внутренняя ошибка сервера"})
		}
		return
	}

	// 3. Успешный ответ
	writeJSON(w, http.StatusOK, map[string]interface{}{})
}
