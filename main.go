package main

import (
	"log"

	"github.com/soundtrack11/go_final_project/pkg/db"
	"github.com/soundtrack11/go_final_project/pkg/server"
)

func main() {
	// Инициализация БД
	if err := db.Init(); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	// Закрытие БД при завершении работы
	defer db.DB.Close()

	// Запуск сервера
	if err := server.Run(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
