package api

import (
	"net/http"
	"time"

	scheduler "github.com/AleksandriUrtaev/go_final_project2/pkg/scheduler"
)

const dateLayout = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {

	// параметры
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateLayout, nowStr)
		if err != nil {
			http.Error(w, "invalid now", http.StatusBadRequest)
			return
		}
	}

	//
	next, err := scheduler.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// ответ
	w.Write([]byte(next))
}
