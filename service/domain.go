package service

import (
	"context"
	"time"
)

// Service ...
type Service struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Versions    int       `json:"versions" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type ServiceRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// Version ...
type Version struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type VersionRequest struct {
	Name string `json:"name"`
}

// ServiceInterface
type ServiceInterface interface {
	Fetch(ctx context.Context, searchBy string, sortBy string, pageNumber int, itemsPerPage int) (res []*Service, err error)
	FetchByID(c context.Context, id int) (res *Service, err error)
	Create(c context.Context, attributes ServiceRequest) (res *Service, err error)
	Update(c context.Context, serviceID int, r ServiceRequest) (res *Service, err error)
	Delete(c context.Context, serviceID int) (err error)
	FetchVersions(c context.Context, serviceID int) (res []*Version, err error)
	CreateVersion(c context.Context, serviceID int, r VersionRequest) (res *Service, err error)
}
