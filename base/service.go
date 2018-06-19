package base

import (
	"log"

	"github.com/syndicatedb/vodka/repositories"
)

// Service - base service interface that implements basic CRUD methods for base controller
type Service interface {
	Find(map[string]interface{}, map[string]interface{}) (interface{}, error)
	FindByID(interface{}) (interface{}, error)
	Create(interface{}) (interface{}, error)
	Save(map[string]interface{}, map[string]interface{}) (interface{}, error)
	Update(map[string]interface{}, map[string]interface{}) (interface{}, error)
	DeleteByID(interface{}) (interface{}, error)
}

type service struct {
	repository repositories.Recorder
}

// NewService - service constructor
func NewService(repo repositories.Recorder) Service {
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

func (s *service) Save(query, payload map[string]interface{}) (interface{}, error) {
	log.Println("BODY", payload)
	log.Println("QUERY", query)
	return s.repository.Save(payload, query)
}

func (s *service) Update(query, payload map[string]interface{}) (interface{}, error) {
	return s.repository.Update(query, payload)
}

func (s *service) DeleteByID(id interface{}) (interface{}, error) {
	return s.repository.DeleteByID(id)
}
