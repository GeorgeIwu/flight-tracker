package service

import (
	"context"
	"fmt"
	"service-catalog/ent"
	"service-catalog/ent/service"
)

type ServiceUsecase struct {
	repo         ServiceRepoInterface
}

// NewServiceUsecase will create new an ServiceUsecase object representation of ServiceInterface
func NewServiceUsecase(r ServiceRepoInterface) *ServiceUsecase {
	return &ServiceUsecase{
		repo:         r,
	}
}

func (u *ServiceUsecase) Fetch(ctx context.Context, searchBy string, sortBy string, pageNumber int, itemsPerPage int) (res []*Service, err error) {
	pageOffset := (pageNumber - 1) * itemsPerPage
	serviceEntities, err := u.repo.Get(ctx, searchBy, sortBy, pageOffset, itemsPerPage)
	if err != nil {
		return nil, err
	}

	services := []*Service{}
	for _, serviceEntity := range serviceEntities {
		newservice, _ := mapService(ctx, serviceEntity)
		services = append(services, newservice)
	}

	return services, nil
}

func (u *ServiceUsecase) FetchByID(ctx context.Context, id int) (res *Service, err error) {

	serviceEntity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	service, err := mapService(ctx, serviceEntity)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (u *ServiceUsecase) Create(c context.Context, r ServiceRequest) (res *Service, err error) {
	serviceEntity, err := u.repo.Create(c, r)
	if err != nil {
		return nil, err
	}

	service, err := mapService(c, serviceEntity)
	if err != nil {
		return nil, err
	}
	fmt.Println("service was created: ", service)

	return service, nil
}

func (u *ServiceUsecase) Update(c context.Context, serviceID int, r ServiceRequest) (res *Service, err error) {

	// Update the Service.
	serviceEntity, err := u.repo.Update(c, serviceID, r)
	if err != nil {
		return nil, err
	}

	service, err := mapService(c, serviceEntity)
	if err != nil {
		return nil, err
	}
	fmt.Println("version was created: ", service)

	return service, nil
}

func (u *ServiceUsecase) Remove(c context.Context, serviceID int) (err error) {

	newerr := u.repo.Delete(c, serviceID)
	if newerr != nil {
		return newerr
	}

	fmt.Println("version was deleted ")

	return nil
}

func (u *ServiceUsecase) FetchVersions(c context.Context, serviceID int) (res []*Version, err error) {

	versionEntities, err := u.repo.GetVersions(c, serviceID)
	if versionEntities == nil {
		return nil, err
	}

	versions := []*Version{}
	for _, versionEntity := range versionEntities {
		newversion := &Version{
			ID:        int64(versionEntity.ID),
			Name:      versionEntity.Name,
			UpdatedAt: versionEntity.UpdatedAt,
		}
		versions = append(versions, newversion)
	}

	return versions, nil
}

func (u *ServiceUsecase) CreateVersion(c context.Context, serviceID int, r VersionRequest) (res *Service, err error) {

	serviceEntity, err := u.repo.CreateVersion(c, serviceID, r)
	if err != nil {
		return nil, err
	}

	service, err := mapService(c, serviceEntity)
	if err != nil {
		return nil, err
	}
	fmt.Println("version was created: ", service)

	return service, nil
}

func mapService(c context.Context, data *ent.Service) (*Service, error) {

	newservice := &Service{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Versions:    data.VersionCount,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return newservice, nil
}

func getSortType(sort string) string {

	switch sort {
	case "title":
		return service.FieldTitle
	case "description":
		return service.FieldDescription
	default:
		return service.FieldID
	}
}

