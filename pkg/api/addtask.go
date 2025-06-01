package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/soundtrack11/go_final_project/pkg/db"
	"github.com/soundtrack11/go_final_project/pkg/nextdate"
)

// Cтруктура входящего запроса
type addTaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// Cтруктура ответа
type addTaskResponse struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// Обработка запросов на добавление задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON запрос
	var req addTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	// Валидация обязательных полей
	if req.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Обработка даты
	now := time.Now().UTC()
	today := now.Format("20060102")

	if req.Date == "" {
		req.Date = today
	}

	// Проверка формата даты
	date, err := time.Parse("20060102", req.Date)
	if err != nil {
		writeJSONError(w, "Неверный формат даты", http.StatusBadRequest)
		return
	}

	// Приводим дату к началу дня для корректного сравнения
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Проверка и коррекция даты
	if date.Before(nowDate) {
		if req.Repeat == "" {
			// Если дата в прошлом и нет повтора - используем сегодня
			req.Date = today
		} else {
			// Если есть повтор - вычисляем следующую дату
			next, err := nextdate.NextDate(nowDate, req.Date, req.Repeat)
			if err != nil {
				writeJSONError(w, "Некорректное правило повторения: "+err.Error(), http.StatusBadRequest)
				return
			}
			req.Date = next
		}
	} else if req.Repeat != "" {
		// Проверяем правило повторения для будущих дат
		_, err := nextdate.NextDate(nowDate, req.Date, req.Repeat)
		if err != nil {
			writeJSONError(w, "Некорректное правило повторения: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Создаем задачу для БД
	task := &db.Task{
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	// Добавляем задачу в БД
	id, err := db.AddTask(task)
	if err != nil {
		writeJSONError(w, "Ошибка добавления задачи: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	writeJSONResponse(w, addTaskResponse{ID: id}, http.StatusOK)
}

// Отправление JSON ответа
func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// Отправление JSON с ошибкой
func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	writeJSONResponse(w, struct {
		Error string `json:"error"`
	}{Error: errorMsg}, statusCode)
}
