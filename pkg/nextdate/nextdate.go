package nextdate

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Вычисляем следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Парсим исходную дату
	startDate, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", errors.New("неверный формат даты")
	}

	// Проверяем наличие правила повторения
	if repeat == "" {
		return "", errors.New("пустое правило повторения")
	}

	// Удаляем лишние пробелы
	repeat = strings.TrimSpace(repeat)
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("пустое правило повторения")
	}

	// Обрабатываем разные типы правил
	switch parts[0] {
	case "d":
		return handleDailyRule(now, startDate, parts)
	case "y":
		return handleYearlyRule(now, startDate)
	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

// Обрабатываем ежедневные правила (d <дни>)
func handleDailyRule(now, startDate time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("неверный формат для правила 'd'")
	}

	// Парсим количество дней
	days, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("неверное число дней")
	}

	// Проверяем допустимость интервала
	if days < 1 || days > 400 {
		return "", errors.New("интервал дней должен быть от 1 до 400")
	}

	// Вычисляем следующую дату
	current := startDate
	for {
		current = current.AddDate(0, 0, days)
		if current.After(now) {
			break
		}
	}
	return current.Format("20060102"), nil
}

// Обрабатываем ежегодные правила (y)
func handleYearlyRule(now, startDate time.Time) (string, error) {
	// Вычисляем следующую дату
	current := startDate
	for {
		// Добавляем 1 год
		nextYear := current.Year() + 1
		nextDate := time.Date(nextYear, current.Month(), current.Day(), 0, 0, 0, 0, time.UTC)

		// Корректируем 29 февраля для невисокосных лет
		if current.Month() == time.February && current.Day() == 29 {
			// Проверяем, существует ли 29 февраля в следующем году
			feb29 := time.Date(nextYear, time.February, 29, 0, 0, 0, 0, time.UTC)
			if feb29.Month() != time.February || feb29.Day() != 29 {
				// Если 29 февраля не существует, используем 1 марта
				nextDate = time.Date(nextYear, time.March, 1, 0, 0, 0, 0, time.UTC)
			}
		}

		current = nextDate

		// Проверяем, что дата после текущего времени
		if current.After(now) {
			break
		}
	}
	return current.Format("20060102"), nil
}
