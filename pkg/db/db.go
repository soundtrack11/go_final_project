package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// Глобальная переменная для доступа к БД
var DB *sql.DB

// Инициализация БД
func Init() error {
	dbFile := getDBPath()

	// Проверяем существование файла БД
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открываем соединение с БД
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbFile)
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Создаем таблицу при необходимости
	if install {
		if err := createSchema(); err != nil {
			return fmt.Errorf("ошибка создания схемы: %w", err)
		}
	}

	// Проверяем соединение
	return DB.Ping()
}

// Получаем путь к файлу БД
func getDBPath() string {
	// 1. Проверяем переменную окружения
	if envPath := os.Getenv("TODO_DBFILE"); envPath != "" {
		return envPath
	}

	// 2. Используем путь по умолчанию
	return "scheduler.db"
}

// Создаем таблицы и индексы
func createSchema() error {
	schema := `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT '',
		title VARCHAR(255) NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT '',
		repeat VARCHAR(128) NOT NULL DEFAULT ''
	);
	
	CREATE INDEX idx_date ON scheduler(date);
	`

	_, err := DB.Exec(schema)
	return err
}
