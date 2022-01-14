package service

import (
	"main/domain"
	"main/repository"
	"net/http"
)

type CarRestService interface {
	Fetch(id string) (http.Response, error)
	Create(domain.CarDTO) (http.Response, error)
	Delete(id, version string) (http.Response, error)
}

type DefaultCarRestService struct {
	repo repository.CarRestRepository
}

func (as DefaultCarRestService) Fetch(id string) (http.Response, error) {
	return as.repo.Fetch(id)
}

func (as DefaultCarRestService) Create(accData domain.CarDTO) (http.Response, error) {
	return as.repo.Create(accData)
}

func (as DefaultCarRestService) Delete(id, version string) (http.Response, error) {
	return as.repo.Delete(id, version)
}

func NewCarRestService(repo repository.CarRestRepository) DefaultCarRestService {
	return DefaultCarRestService{repo}
}
