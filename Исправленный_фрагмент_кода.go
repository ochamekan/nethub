package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

// ВНИМАНИЕ: в этом фрагменте есть несколько ошибок и плохих практик.
// Кандидату нужно:
// 1) Найти и описать проблемы.
// 2) Предложить, как переписать код лучше.

var db *sql.DB

func initDB() error {
	// Потенциальная проблема: ошибка игнорируется
	// ИСПРАВЛЕНО: обрабатываем ошибку sql.Open
	db, err := sql.Open("postgres", "host=localhost user=app dbname=devices sslmode=disable")
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}
	defer db.Close()

	// Нет проверки доступности соединения и таймаута
	// ИСПРАВЛЕНО: проверяем реальную доступность БД с таймаутом, пример взял с гитхаба
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("db.Ping: %w", err)
	}
	return nil
}

// Device простая модель устройства
type Device struct {
	ID       int64
	Hostname string
	IP       string
}

// deviceHandler получает устройство по id и пишет в лог таблицу audit_log
func deviceHandler(w http.ResponseWriter, r *http.Request) {
	// ИСПРАВЛЕНО: контекст с таймаутом для всего запроса
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	// Потенциальная проблема: контекст без таймаута, возможная утечка
	// ИСПРАВЛЕНО: Избавляемся от утечки
	go func() {
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("long debug operation finished")
		case <-ctx.Done():
			fmt.Println("debug operation cancelled:", ctx.Err())
		}
	}()

	// Потенциальная проблема: строка запроса строится конкатенацией
	// ИСПРАВЛЕНО: Неправильное использование ф-ии, аргументы передаем в параметры QueryRowContext, чтобы избежать sql-injection
	row := db.QueryRowContext(ctx, "SELECT id, hostname, ip FROM devices WHERE id = $1", idStr)

	var d Device
	err := row.Scan(&d.ID, &d.Hostname, &d.IP)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	// Потенциальная проблема: игнорируется ошибка вставки в audit_log
	// ИСПРАВЛЕНО: Обрабатываем ошибку
	_, err = db.ExecContext(ctx, "INSERT INTO audit_log(device_id, ts, action) VALUES ($1, now(), 'view')", d.ID)
	if err != nil {
		// Я бы отправил лог ошибки, но не возвращал бы http.Error, т.к. если проблема с бд логов, то мы вернем 500 при корректном запросе
		// Как минимум
		log.Printf("failed to save log about device %s: %s\n", d.ID, err.Error())
	}

	// Потенциальная проблема: нет установки Content-Type
	// ИСПРАВЛЕНО: Сетим тип контента в хедер ответа
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("Device: %s (%s)", d.Hostname, d.IP)))
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/device", deviceHandler)

	// Потенциальная проблема: сервер никогда не завершится, ошибки ListenAndServe игнорируются
	// ИСПРАВЛЕНО: Обрабатываем ошибку
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
