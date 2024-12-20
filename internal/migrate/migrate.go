package migrate

import (
	"database/sql"
	"fmt"

	"github.com/Alexander272/Pinger/internal/migrate/db/migrations"
	"github.com/Alexander272/Pinger/pkg/logger"
	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(&migrations.Content)

	// dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
	// 	conf.Username, conf.Password, conf.Host, conf.Port, conf.DbName, conf.SSLMode,
	// )

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// db, err := goose.OpenDBWithDriver("pgx", dsn)
	// if err != nil {
	// 	return fmt.Errorf("failed to open db: %w", err)
	// }

	logger.Info("migration up till last")
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to migrate up: %w", err)
	}

	// err = db.Close()
	// if err != nil {
	// 	return fmt.Errorf("failed to close db: %w", err)
	// }
	return nil
}
