package service

import (
	"context"
	"fmt"
	"servicecatalog/ent"
	"servicecatalog/ent/predicate"
	"servicecatalog/ent/service"
	"time"
)

type ServiceUsecase struct {
	client         *ent.Client
	contextTimeout time.Duration
}

// NewServiceUsecase will create new an ServiceUsecase object representation of ServiceInterface
func NewServiceUsecase(c *ent.Client, timeout time.Duration) *ServiceUsecase {
	return &ServiceUsecase{
		client:         c,
		contextTimeout: timeout,
	}
}

func (u *ServiceUsecase) Fetch(ctx context.Context, searchBy string, sortBy string, pageNumber int, itemsPerPage int) (res []*Service, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	pageOffset := (pageNumber - 1) * itemsPerPage
	serviceEntities, err := u.client.Service.Query().
		Where(
			service.Or(
				service.TitleContains(searchBy),
				service.DescriptionContains(searchBy),
			),
		).
		Offset(pageOffset).
		Limit(itemsPerPage).
		Order(ent.Desc(getSortType(sortBy))).
		All(ctx)
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

func (u *ServiceUsecase) FetchByID(c context.Context, id int) (res *Service, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	serviceEntity, err := u.client.Service.Get(ctx, id)
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
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	tx, err := u.client.Tx(ctx)
	defer cancel()

	// Create a version.
	version, err := tx.Version.Create().Save(ctx)
	if err != nil {
		return nil, rollback(tx, err)
	}
	fmt.Println("version was created: ", version)

	// Create a new Service.
	serviceEntity, err := tx.Service.
		Create().
		SetTitle(r.Title).
		SetDescription(r.Description).
		SetVersionCount(1).
		SetUpdatedAt(time.Now()).
		SetCreatedAt(time.Now()).
		AddVersions(version).
		Save(ctx)
	if err != nil {
		return nil, rollback(tx, err)
	}

	service, err := mapService(ctx, serviceEntity)
	if err != nil {
		return nil, rollback(tx, err)
	}
	tx.Commit()
	fmt.Println("service was created: ", service)

	return service, nil
}

func (u *ServiceUsecase) Update(c context.Context, serviceID int, r ServiceRequest) (res *Service, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	tx, err := u.client.Tx(ctx)
	defer cancel()

	existingService, err := tx.Service.Get(ctx, serviceID)
	if existingService == nil {
		return nil, err
	}

	// Update the Service.
	serviceEntity, err := tx.Service.UpdateOneID(serviceID).
		SetTitle(r.Title).
		SetDescription(r.Description).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, rollback(tx, err)
	}

	service, err := mapService(ctx, serviceEntity)
	if err != nil {
		return nil, rollback(tx, err)
	}
	tx.Commit()
	fmt.Println("version was created: ", service)

	return service, nil
}

func (u *ServiceUsecase) Delete(c context.Context, serviceID int) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	existingService, err := u.client.Service.Get(ctx, serviceID)
	if existingService == nil {
		return err
	}

	// Delete the Service.
	newerr := u.client.Service.DeleteOneID(serviceID).Exec(ctx)
	if newerr != nil {
		return newerr
	}

	fmt.Println("version was deleted ")

	return nil
}

func (u *ServiceUsecase) FetchVersions(c context.Context, serviceID int) (res []*Version, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	versionEntities, err := u.client.Version.Query().
		Where(
			predicate.Version(service.ID(serviceID)),
		).All(ctx)
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
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	tx, err := u.client.Tx(ctx)
	defer cancel()

	existingService, err := tx.Service.Get(ctx, serviceID)
	if existingService == nil {
		return nil, err
	}

	// Create a version.
	version, err := tx.Version.Create().SetName(r.Name).Save(ctx)
	if err != nil {
		return nil, rollback(tx, err)
	}
	fmt.Println("version was created: ", version)

	// Update the Service.
	serviceEntity, err := tx.Service.UpdateOneID(serviceID).SetVersionCount(existingService.VersionCount + 1).Save(ctx)
	if err != nil {
		return nil, rollback(tx, err)
	}

	service, err := mapService(ctx, serviceEntity)
	if err != nil {
		return nil, rollback(tx, err)
	}
	tx.Commit()
	fmt.Println("version was created: ", service)

	return service, nil
}

func mapService(c context.Context, data *ent.Service) (*Service, error) {

	newservice := &Service{
		ID:          int64(data.ID),
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

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}
