package repository

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
)

type Repository interface {
	Migrate() error

	GetUserByID(ctx context.Context, id uint) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	CheckUsernameExists(ctx context.Context, username string) (bool, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CreateUser(ctx context.Context, user *model.User) error
}
