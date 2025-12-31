package sql

import (
	"context"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"gorm.io/gorm"
)

func (r *sqlRepository) GetUserByID(
	ctx context.Context,
	id uint,
) (model.User, error) {
	user, err := gorm.
		G[model.User](r.db).
		Where("id = ?", id).
		First(ctx)
	return user, err
}

func (r *sqlRepository) GetUserByUsername(
	ctx context.Context,
	username string,
) (model.User, error) {
	user, err := gorm.
		G[model.User](r.db).
		Where("username = ?", username).
		First(ctx)
	return user, err
}

func (r *sqlRepository) UpdateUser(
	ctx context.Context,
	user *model.User,
) error {
	_, err := gorm.
		G[model.User](r.db).
		Where("id = ?", user.ID).
		Updates(ctx, model.User{AvatarURL: user.AvatarURL, Bio: user.Bio})
	return err
}

func (r *sqlRepository) CheckUsernameExists(
	ctx context.Context,
	username string,
) (bool, error) {
	users, err := gorm.
		G[model.User](r.db).
		Where("username = ?", username).
		Limit(1).
		Find(ctx)
	if err != nil {
		return false, err
	}
	return len(users) > 0, nil
}

func (r *sqlRepository) CheckEmailExists(
	ctx context.Context,
	email string,
) (bool, error) {
	users, err := gorm.
		G[model.User](r.db).
		Where("email = ?", email).
		Limit(1).
		Find(ctx)
	if err != nil {
		return false, err
	}
	return len(users) > 0, nil
}

func (r *sqlRepository) CreateUser(
	ctx context.Context,
	user *model.User,
) error {
	err := gorm.
		G[model.User](r.db).
		Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
