package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/soundtrack11/go_final_project/pkg/db"
	"github.com/soundtrack11/go_final_project/pkg/nextdate"
)

// Структура для ответа с задачей
type taskResponse struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Обработка запросов на получение задачи по ID
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	// Формируем ответ
	resp := taskResponse{
		ID:      strconv.FormatInt(task.ID, 10),
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeat:  task.Repeat,
	}

	writeJSONResponse(w, resp, http.StatusOK)
}

// Структура для запроса обновления задачи
type updateTaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Обработка запросов на обновление задачи
func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Декодируем JSON запрос
	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Ошибка декодирования JSON", http.StatusBadRequest)
		return
	}

	// Проверяем обязательные поля
	if req.Title == "" {
		writeJSONError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Преобразуем ID в int64
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		writeJSONError(w, "Некорректный идентификатор", http.StatusBadRequest)
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
		ID:      id,
		Date:    req.Date,
		Title:   req.Title,
		Comment: req.Comment,
		Repeat:  req.Repeat,
	}

	// Обновляем задачу в БД
	if err := db.UpdateTask(task); err != nil {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем успешный ответ
	writeJSONResponse(w, struct{}{}, http.StatusOK)
}
