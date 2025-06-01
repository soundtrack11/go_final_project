package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/soundtrack11/go_final_project/pkg/api"
)

// Run запускает веб-сервер
func Run() error {
	// Создаем собственный мультиплексор
	mux := http.NewServeMux()

	// Регистрируем API обработчики
	api.Init(mux)

	// Регистрируем файловый сервер
	mux.Handle("/", http.FileServer(http.Dir("web")))

	// Запускаем сервер
	port := getPort()
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

// Определяем порт для сервера
func getPort() int {
	// 1. Проверяем переменную окружения TODO_PORT
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			return p
		}
	}

	// 2. Используем порт по умолчанию
	return 7540
}
