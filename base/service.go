package base

import (
	"github.com/niklucky/vodka/repositories"
)

type Service interface {
	Find(map[string]interface{}, map[string]interface{}) (interface{}, error)
	FindByID(interface{}) (interface{}, error)
	Create(interface{}) (interface{}, error)
	Update(map[string]interface{}, map[string]interface{}) (interface{}, error)
	DeleteByID(interface{}) (interface{}, error)
}

type service struct {
	repository repositories.Recorder
}

func NewService(repo repositories.Recorder) *service {
	return &service{
		repository: repo,
	}
}

func (s *service) FindByID(id interface{}) (interface{}, error) {
	return s.repository.FindByID(id)
}

func (s *service) Find(query, params map[string]interface{}) (interface{}, error) {
	return s.repository.Find(query, params)
}

func (s *service) Create(payload interface{}) (interface{}, error) {
	return s.repository.Create(payload)
}

func (s *service) Update(query, payload map[string]interface{}) (interface{}, error) {
	return s.repository.Update(query, payload)
}

func (s *service) DeleteByID(id interface{}) (interface{}, error) {
	return s.repository.DeleteByID(id)
}
