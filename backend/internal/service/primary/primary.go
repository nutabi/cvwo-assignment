package primary

import (
	"github.com/nutabi/cvwo-assignment/backend/internal/repository"
	"github.com/nutabi/cvwo-assignment/backend/internal/service"
)

type primaryService struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) service.Service {
	return &primaryService{repo: repo}
}

func (s *primaryService) Repo() repository.Repository {
	return s.repo
}
