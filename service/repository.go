package service

import (
	"context"
	"fmt"
	"service-catalog/ent"
	"service-catalog/ent/predicate"
	"service-catalog/ent/service"
	"time"
)

type ServiceRepo struct {
	client         *ent.Client
	contextTimeout time.Duration
}

// NewServiceRepo will create new an ServiceRepo object representation of ServiceInterface
func NewServiceRepo(c *ent.Client, timeout time.Duration) *ServiceRepo {
	return &ServiceRepo{
		client:         c,
		contextTimeout: timeout,
	}
}

func (u *ServiceRepo) Get(ctx context.Context, searchBy string, sortBy string, pageOffset int, itemsPerPage int) (res []*ent.Service, err error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

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

	return serviceEntities, nil
}

func (u *ServiceRepo) GetByID(c context.Context, id int) (res *ent.Service, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	serviceEntity, err := u.client.Service.Get(ctx, id)
	if err != nil {
		return nil, err
	}


	return serviceEntity, nil
}

func (u *ServiceRepo) Create(c context.Context, r ServiceRequest) (res *ent.Service, err error) {
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

	tx.Commit()
	fmt.Println("service was created: ", serviceEntity)

	return serviceEntity, nil
}

func (u *ServiceRepo) Update(c context.Context, serviceID int, r ServiceRequest) (res *ent.Service, err error) {
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

	tx.Commit()
	fmt.Println("version was created: ", serviceEntity)

	return serviceEntity, nil
}

func (u *ServiceRepo) Delete(c context.Context, serviceID int) (err error) {
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

func (u *ServiceRepo) GetVersions(c context.Context, serviceID int) (res []*ent.Version, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	versionEntities, err := u.client.Version.Query().
		Where(
			predicate.Version(service.ID(serviceID)),
		).All(ctx)
	if versionEntities == nil {
		return nil, err
	}

	return versionEntities, nil
}

func (u *ServiceRepo) CreateVersion(c context.Context, serviceID int, r VersionRequest) (res *ent.Service, err error) {
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

	tx.Commit()
	fmt.Println("version was created: ", serviceEntity)

	return serviceEntity, nil
}

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}
