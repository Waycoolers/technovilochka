package storage

import (
	"technovilochka/internal/platform/config"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestBuildDSNEscapesCredentials(t *testing.T) {
	cfg := &config.DBConfig{
		Host: "db",
		Port: "5432",
		User: "bot:user",
		Pass: "p@ss:w/rd?#%+ x",
		Name: "technovilochka",
	}

	parsed, err := pgconn.ParseConfig(buildDSN(cfg))
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}

	if parsed.User != cfg.User || parsed.Password != cfg.Pass || parsed.Host != cfg.Host ||
		parsed.Port != 5432 || parsed.Database != cfg.Name {
		t.Errorf("dsn round trip mismatch: user=%q pass=%q host=%q port=%d db=%q",
			parsed.User, parsed.Password, parsed.Host, parsed.Port, parsed.Database)
	}
}
