package storage

import (
	"log/slog"
	"net"
	"net/url"
	"technovilochka/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib" // регистрирует драйвер database/sql "pgx"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(cfg *config.DBConfig) (*Storage, error) {
	db, err := sqlx.Connect("pgx", buildDSN(cfg))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	slog.Info("Database connection successfully")

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() {
	err := s.db.Close()
	if err != nil {
		slog.Error("Error closing database", "error", err)
	}
}

// buildDSN собирает URL подключения через url.URL, чтобы логин и пароль экранировались:
// при сборке через Sprintf символы @ : / ? # в пароле ломают разбор DSN.
func buildDSN(cfg *config.DBConfig) string {
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Pass),
		Host:     net.JoinHostPort(cfg.Host, cfg.Port),
		Path:     "/" + cfg.Name,
		RawQuery: "sslmode=disable",
	}
	return dsn.String()
}
