package sql

import (
	"log/slog"

	"github.com/nutabi/cvwo-assignment/backend/internal/model"
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"gorm.io/gorm"
)

type sqlRepository struct {
	db *gorm.DB
}

func Connect(dialect gorm.Dialector) (repository.Repository, error) {
	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		slog.Error("unable to connect to SQL database", "error", err)
		return nil, err
	}
	slog.Debug("connected to SQL database")
	return &sqlRepository{db: db}, nil
}

func (r *sqlRepository) Migrate() error {
	slog.Debug("running database migrations")

	// Migrate user
	err := r.db.AutoMigrate(&model.User{})
	if err != nil {
		slog.Error("failed to migrate User model", "error", err)
		return err
	}

	// Migrate topic
	err = r.db.AutoMigrate(&model.Topic{})
	if err != nil {
		slog.Error("failed to migrate Topic model", "error", err)
		return err
	}

	// Migrate post
	err = r.db.AutoMigrate(&model.Post{})
	if err != nil {
		slog.Error("failed to migrate Post model", "error", err)
		return err
	}

	// Migrate comment
	err = r.db.AutoMigrate(&model.Comment{})
	if err != nil {
		slog.Error("failed to migrate Comment model", "error", err)
		return err
	}

	slog.Debug("database migrations completed")

	return nil
}
