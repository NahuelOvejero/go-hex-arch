package domain

import (
	"net/http"
)

type CarRepositoryPort interface {
	Fetch(id string) (http.Response, error)
	Create(CarDTO) (http.Response, error)
	Delete(id string, version string) (http.Response, error)
}

type CarData struct {
	ID      string `json:"id,omitempty"`
	Owner   *User  `json:"Owner,omitempty"`
	Type    string `json:"type,omitempty"`
	Version *int64 `json:"version,omitempty"`
}

type User struct {
	Id      int64  `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

type CarDTO struct {
	Data *CarData
}
