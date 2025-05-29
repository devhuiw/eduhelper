package mysql

import (
	"database/sql"
	"fmt"
	"service/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

// New открывает соединение с MySQL и возвращает *sql.DB.
func New(cfg config.SQLPath) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.User, cfg.Password, cfg.Host, fmt.Sprintf("%d", cfg.Port), cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Проверка соединения с базой
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}

	return db, nil
}
