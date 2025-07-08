package api

import (
	"net/http"
	"strconv"

	"github.com/soundtrack11/go_final_project/pkg/db"
)

// Структура для ответа со списком задач
type tasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

// Обработка запросов на получение списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим параметр limit
	limit := 50
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	// Получаем задачи из БД
	tasks, err := db.GetTasks(limit)
	if err != nil {
		writeJSONError(w, "Ошибка получения задач: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Формируем ответ
	writeJSONResponse(w, tasksResponse{Tasks: tasks}, http.StatusOK)
}
