package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/soundtrack11/go_final_project/pkg/db"
	"github.com/soundtrack11/go_final_project/pkg/nextdate"
)

// Обработка запросов на отметку задачи выполненной
func doneHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID из параметров запроса
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		writeJSONError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в int64
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		writeJSONError(w, "Некорректный идентификатор", http.StatusBadRequest)
		return
	}

	// Получаем задачу из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обрабатываем задачу в зависимости от типа повторения
	if task.Repeat == "" {
		// Удаляем одноразовую задачу
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		// Вычисляем следующую дату для повторяющейся задачи
		now := time.Now().UTC()
		next, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, "Ошибка вычисления следующей даты: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Обновляем дату выполнения
		if err := db.UpdateTaskDate(id, next); err != nil {
			writeJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Возвращаем успешный ответ
	writeJSONResponse(w, struct{}{}, http.StatusOK)
}
