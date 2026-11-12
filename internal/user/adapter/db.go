package adapter

import (
	"context"
	"technovilochka/internal/user/domain"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, chat_id, username, first_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			chat_id = EXCLUDED.chat_id,
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name;
	`, user.ID, user.ChatID, user.Username, user.FirstName)
	return err
}
