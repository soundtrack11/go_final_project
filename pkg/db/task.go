package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
)

// Структура задачи в системе
type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Сериализация для Task
func (t Task) MarshalJSON() ([]byte, error) {
	// Преобразуем ID в строку
	return json.Marshal(struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}{
		ID:      strconv.FormatInt(t.ID, 10),
		Date:    t.Date,
		Title:   t.Title,
		Comment: t.Comment,
		Repeat:  t.Repeat,
	})
}

// Добавляем новую задачу в базу данных
func AddTask(task *Task) (int64, error) {
	res, err := DB.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Возвращаем список задач, отсортированных по дате
func GetTasks(limit int) ([]Task, error) {
	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		ORDER BY date ASC, id ASC
		LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Возврат пустого массива вместо nil
	if tasks == nil {
		tasks = []Task{}
	}

	return tasks, nil
}

// Возвращаем задачу по ID
func GetTask(id int64) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := DB.QueryRow(query, id)

	var task Task
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("задача не найдена")
		}
		return nil, err
	}
	return &task, nil
}

// Обновляем существующую задачу
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}
	return nil
}

// Удаляем задачу по ID
func DeleteTask(id int64) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}
	return nil
}

// Обновляем дату выполнения задачи
func UpdateTaskDate(id int64, newDate string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, newDate, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("задача не найдена")
	}
	return nil
}
