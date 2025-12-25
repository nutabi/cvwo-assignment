package repository

import (
	"context"
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

type sqlRepository struct {
	db gorm.DB
}

func ConnectSQL(dialect gorm.Dialector) (Repository, error) {
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		slog.Error("Unable to connect to SQL database", "error", err)
		return nil, err
	}
	return &sqlRepository{db: *db}, nil
}

func (r *sqlRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (model.User, error) {
	user, err := gorm.
		G[model.User](&r.db).
		Where("username = ?", username).
		First(ctx)
	return user, err
}
