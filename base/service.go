package base

import (
	"github.com/niklucky/vodka/repositories"
)

type Service struct {
	repository repositories.Recorder
}

func NewService(repo repositories.Recorder) Service {
	return Service{
		repository: repo,
	}
}

func (s *Service) FindByID(id interface{}) (interface{}, error) {
	return s.repository.FindByID(id)
}
