package repository

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
)

type Repository interface {
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
}
