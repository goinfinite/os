package apiController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type VirtualHostController struct {
	virtualHostService *service.VirtualHostService
}

func NewVirtualHostController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *VirtualHostController {
	return &VirtualHostController{
		virtualHostService: service.NewVirtualHostService(persistentDbSvc, trailDbSvc),
	}
}

// ReadVirtualHosts	 godoc
// @Summary      ReadVirtualHosts
// @Description  List virtual hosts.
// @Tags         vhosts
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        hostname query  string  false  "Hostname"
// @Param        type query  string  false  "Type"
// @Param        rootDirectory query  string  false  "RootDirectory"
// @Param        parentHostname query  string  false  "ParentHostname"
// @Param        withMappings query  bool  false  "WithMappings"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadVirtualHostsResponse
// @Router       /v1/vhosts/ [get]
func (controller *VirtualHostController) Read(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.virtualHostService.Read(requestInputData),
	)
}

// CreateVirtualHost    godoc
// @Summary      CreateVirtualHost
// @Description  Create a new vhost.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createVirtualHostDto 	  body    dto.CreateVirtualHost  true  "Only hostname is required.<br />type may be 'top-level', 'subdomain', 'wildcard' or 'alias'. If is not provided, it will be 'top-level'. If type is 'alias', 'parentHostname' will be required."
// @Success      201 {object} object{} "VirtualHostCreated"
// @Router       /v1/vhosts/ [post]
func (controller *VirtualHostController) Create(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.virtualHostService.Create(requestInputData),
	)
}

// DeleteVirtualHost godoc
// @Summary      DeleteVirtualHost
// @Description  Delete a vhost.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        hostname path string true "Hostname to delete"
// @Success      200 {object} object{} "VirtualHostDeleted"
// @Router       /v1/vhosts/{hostname}/ [delete]
func (controller *VirtualHostController) Delete(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.virtualHostService.Delete(requestInputData),
	)
}

// ReadVirtualHostsWithMappings	 godoc
// @Summary      ReadVirtualHostsWithMappings
// @Description  List virtual hosts with mappings.
// @Tags         vhosts
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        hostname query  string  false  "Hostname"
// @Param        type query  string  false  "Type"
// @Param        rootDirectory query  string  false  "RootDirectory"
// @Param        parentHostname query  string  false  "ParentHostname"
// @Param        withMappings query  bool  false  "WithMappings"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.VirtualHostWithMappings
// @Router       /v1/vhosts/mapping/ [get]
func (controller *VirtualHostController) ReadWithMappings(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.virtualHostService.ReadWithMappings(requestInputData),
	)
}

// CreateVirtualHostMapping godoc
// @Summary      CreateVirtualHostMapping
// @Description  Create a new vhost mapping.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createMappingDto	body dto.CreateMapping	true	"hostname, path and targetType are required.<br />matchPattern may be 'begins-with', 'contains', 'equals' or 'ends-with'. If is not provided, it will be 'begins-with'.<br />targetType may be 'url', 'service', 'response-code', 'inline-html' or 'static-files'. If targetType is 'url', targetHttpResponseCode may be provided. If is not provided, targetHttpResponseCode will be '200'. If targetType is 'response-code', targetHttpResponseCode may be provided. If is not provided, targetValue will be required. If both were provided, targetValue will have priority.<br />targetValue must have the same value as the targetType requires."
// @Success      201 {object} object{} "MappingCreated"
// @Router       /v1/vhosts/mapping/ [post]
func (controller *VirtualHostController) CreateMapping(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.virtualHostService.CreateMapping(requestInputData),
	)
}

// DeleteVirtualHostMapping godoc
// @Summary      DeleteVirtualHostMapping
// @Description  Delete a vhost mapping.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        mappingId path uint true "MappingId to delete."
// @Success      200 {object} object{} "MappingDeleted"
// @Router       /v1/vhosts/mapping/{mappingId}/ [delete]
func (controller *VirtualHostController) DeleteMapping(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.virtualHostService.DeleteMapping(requestInputData),
	)
}
