package usecase

import (
	"context"
	"technovilochka/internal/user/domain"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

type RegisterInput struct {
	ID        int64
	ChatID    int64
	Username  string
	FirstName string
}

func (s *Service) Register(ctx context.Context, input RegisterInput) error {
	user := &domain.User{
		ID:        input.ID,
		ChatID:    input.ChatID,
		Username:  input.Username,
		FirstName: input.FirstName,
	}
	return s.repo.Save(ctx, user)
}
