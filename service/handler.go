package service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ServiceHandler  represent the httphandler for Service
type ServiceHandler struct {
	usecase ServiceInterface
}

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	errInternalServerError = "Server Error"
	// ErrNotFound will throw if the requested item is not exists
	errNotFound = "not found"
	// ErrConflict will throw if the current action already exists
	errConflict = "already exist"
	// ErrBadParamInput will throw if the given request-body or params is not valid
	errBadParamInput = "not valid"
)

// NewServiceHandler will initialize the Services/ resources endpoint
func NewServiceHandler(e *echo.Echo, is ServiceInterface) *ServiceHandler {
	handler := &ServiceHandler{
		usecase: is,
	}

	e.GET("/services", handler.Fetch)
	e.GET("/services/:id", handler.FetchByID)
	e.POST("/services", handler.Create)
	e.PUT("/services/:id", handler.Update)
	e.DELETE("/services/:id", handler.Delete)
	e.POST("/services/:id/versions", handler.CreateVersion)
	e.GET("/services/:id/versions", handler.FetchVersions)
	return handler
}

// FetchService will fetch the Service based on given params
func (a *ServiceHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()
	searchBy := c.QueryParam("search")
	sortBy := c.QueryParam("sort")
	page := c.QueryParam("page")
	itemsPerPage := 12

	pageNumber, _ := strconv.Atoi(page)
	if pageNumber == 0 {
		pageNumber = 1
	}

	validate := validator.New()
	err := validate.Var(pageNumber, "gt=0")
	if err != nil {
		return c.JSON(http.StatusNotFound, errBadParamInput)
	}

	services, err := a.usecase.Fetch(ctx, searchBy, sortBy, pageNumber, itemsPerPage)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, services)
}

// FetchByID will get Service by given id
func (a *ServiceHandler) FetchByID(c echo.Context) error {
	ctx := c.Request().Context()
	serviceID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, errNotFound)
	}

	id := int(serviceID)

	service, err := a.usecase.FetchByID(ctx, id)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, service)
}

// Create will create the Service by given request body
func (a *ServiceHandler) Create(c echo.Context) (err error) {
	ctx := c.Request().Context()
	var request ServiceRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(request); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	newservice, err := a.usecase.Create(ctx, request)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, newservice)
}

// Update will update the Service by given request body
func (a *ServiceHandler) Update(c echo.Context) (err error) {
	ctx := c.Request().Context()
	serviceID, err := strconv.Atoi(c.Param("id"))

	var request ServiceRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(request); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	newservice, err := a.usecase.Update(ctx, serviceID, request)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, newservice)
}

// Delete will update the Service by given request body
func (a *ServiceHandler) Delete(c echo.Context) (err error) {
	ctx := c.Request().Context()
	serviceID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	err = a.usecase.Remove(ctx, serviceID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, nil)
}

// FetchVersions will fetch the Versions based on given params
func (a *ServiceHandler) FetchVersions(c echo.Context) error {
	ctx := c.Request().Context()
	serviceID, err := strconv.Atoi(c.Param("id"))

	versions, err := a.usecase.FetchVersions(ctx, serviceID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, versions)
}

// CreateVersion will create the Version by given request body
func (a *ServiceHandler) CreateVersion(c echo.Context) (err error) {
	serviceID, err := strconv.Atoi(c.Param("id"))

	var request VersionRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(request); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	newservice, err := a.usecase.CreateVersion(ctx, serviceID, request)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, newservice)
}

func isRequestValid(m interface{}) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	fmt.Println("error message", err)

	switch {
	case strings.Contains(err.Error(), errInternalServerError):
		return http.StatusInternalServerError
	case strings.Contains(err.Error(), errNotFound):
		return http.StatusNotFound
	case strings.Contains(err.Error(), errConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
