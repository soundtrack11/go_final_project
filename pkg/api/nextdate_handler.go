package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/soundtrack11/go_final_project/pkg/nextdate"
)

const dateLayout = "20060102"

// Обработка запросов для вычисления следующей даты
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Извлекаем параметры из запроса
	nowParam := r.URL.Query().Get("now")
	dateParam := r.URL.Query().Get("date")
	repeatRule := r.URL.Query().Get("repeat")

	// Обрабатываем параметр now
	nowTime := time.Now().UTC()
	if nowParam != "" {
		parsedNow, err := time.Parse(dateLayout, nowParam)
		if err != nil {
			http.Error(w, "Неверный формат параметра now", http.StatusBadRequest)
			return
		}
		nowTime = parsedNow
	}

	// Проверяем обязательные параметры
	if dateParam == "" {
		http.Error(w, "Параметр date обязателен", http.StatusBadRequest)
		return
	}
	if repeatRule == "" {
		http.Error(w, "Параметр repeat обязателен", http.StatusBadRequest)
		return
	}

	// Вычисляем следующую дату
	nextDate, err := nextdate.NextDate(nowTime, dateParam, repeatRule)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, nextDate)
}
