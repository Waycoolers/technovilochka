package usecase

import (
	"context"
	"technovilochka/internal/user/domain"
)

type Repository interface {
	Save(context.Context, *domain.User) error
}
