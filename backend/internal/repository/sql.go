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

func (r *sqlRepository) GetUserByID(
	ctx context.Context,
	id uint,
) (model.User, error) {
	user, err := gorm.
		G[model.User](&r.db).
		Where("id = ?", id).
		First(ctx)
	return user, err
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

func (r *sqlRepository) UpdateUser(
	ctx context.Context,
	user *model.User,
) error {
	_, err := gorm.
		G[model.User](&r.db).
		Where("id = ?", user.ID).
		Updates(ctx, model.User{AvatarURL: user.AvatarURL, Bio: user.Bio})
	return err
}
