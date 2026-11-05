package config

import (
	"errors"
	"os"
)

type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func loadDBConfig() (*DBConfig, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || pass == "" || name == "" {
		return nil, errors.New("environments for database configuration are not set")
	}

	return &DBConfig{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
		Name: name,
	}, nil
}
