package repository

import (
	"context"
	"errors"
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

func (r *sqlRepository) Migrate() error {
	// Migrate user
	err := r.db.AutoMigrate(&model.User{})
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
		return err
	}

	// Migrate topic
	err = r.db.AutoMigrate(&model.Topic{})
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
		return err
	}

	// Migrate post
	err = r.db.AutoMigrate(&model.Post{})
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
		return err
	}

	// Migrate comment
	err = r.db.AutoMigrate(&model.Comment{})
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
		return err
	}

	return nil
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

func (r *sqlRepository) CheckUsernameExists(
	ctx context.Context,
	username string,
) (bool, error) {
	_, err := gorm.
		G[model.User](&r.db).
		Where("username = ?", username).
		Limit(1).
		Find(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *sqlRepository) CheckEmailExists(
	ctx context.Context,
	email string,
) (bool, error) {
	_, err := gorm.
		G[model.User](&r.db).
		Where("email = ?", email).
		First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *sqlRepository) CreateUser(
	ctx context.Context,
	user *model.User,
) error {
	err := gorm.
		G[model.User](&r.db).
		Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
