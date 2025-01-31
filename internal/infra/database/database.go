package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/VieiraVitor/transaction-flow/config"
	_ "github.com/lib/pq"
)

func ConnectDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	var db *sql.DB
	var err error

	for i := 1; i <= 5; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}

		fmt.Printf("Attempt %d: Failed to connect to database: %v\n", i, err)
		time.Sleep(time.Duration(i) * time.Second) // Exponencial Backoff
	}

	return nil, fmt.Errorf("failed to connect to database after 5 attempts: %v", err)
}
