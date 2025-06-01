package api

import (
	"net/http"
	"strconv"

	"github.com/soundtrack11/go_final_project/pkg/db"
)

// Регистрация обработчиков API
func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", NextDateHandler)
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
	mux.HandleFunc("/api/task/done", doneHandler)
}

// Обработка запросов к /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработка запросов на удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	// Удаляем задачу
	if err := db.DeleteTask(id); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем успешный ответ
	writeJSONResponse(w, struct{}{}, http.StatusOK)
}
