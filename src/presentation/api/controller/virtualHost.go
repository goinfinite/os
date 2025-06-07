package apiController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/labstack/echo/v4"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
)

type VirtualHostController struct {
	virtualHostLiaison *liaison.VirtualHostLiaison
}

func NewVirtualHostController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *VirtualHostController {
	return &VirtualHostController{
		virtualHostLiaison: liaison.NewVirtualHostLiaison(persistentDbSvc, trailDbSvc),
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
// @Router       /v1/vhost/ [get]
func (controller *VirtualHostController) Read(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.Read(requestInputData),
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
// @Router       /v1/vhost/ [post]
func (controller *VirtualHostController) Create(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.Create(requestInputData),
	)
}

// UpdateVirtualHost    godoc
// @Summary      UpdateVirtualHost
// @Description  Update a vhost.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateVirtualHostDto 	  body    dto.UpdateVirtualHost  true  "Only hostname is required."
// @Success      200 {object} object{} "VirtualHostUpdated"
// @Router       /v1/vhost/ [put]
func (controller *VirtualHostController) Update(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.Update(requestInputData),
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
// @Router       /v1/vhost/{hostname}/ [delete]
func (controller *VirtualHostController) Delete(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.Delete(requestInputData),
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
// @Param        createdBeforeAt query  integer  false  "CreatedBeforeAt (Unix timestamp)"
// @Param        createdAfterAt query  integer  false  "CreatedAfterAt (Unix timestamp)"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.VirtualHostWithMappings
// @Router       /v1/vhost/mapping/ [get]
func (controller *VirtualHostController) ReadWithMappings(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.ReadWithMappings(requestInputData),
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
// @Router       /v1/vhost/mapping/ [post]
func (controller *VirtualHostController) CreateMapping(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.CreateMapping(requestInputData),
	)
}

// UpdateVirtualHostMapping godoc
// @Summary      UpdateVirtualHostMapping
// @Description  Update a vhost mapping.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateMappingDto body dto.UpdateMapping true "Only id is required. Other fields are optional and will only be updated if provided."
// @Success      200 {object} object{} "MappingUpdated"
// @Router       /v1/vhost/mapping/ [put]
func (controller *VirtualHostController) UpdateMapping(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.UpdateMapping(requestInputData),
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
// @Router       /v1/vhost/mapping/{mappingId}/ [delete]
func (controller *VirtualHostController) DeleteMapping(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.DeleteMapping(requestInputData),
	)
}

// ReadMappingSecurityRules godoc
// @Summary      ReadMappingSecurityRules
// @Description  List mapping security rules.
// @Tags         vhosts
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id query  string  false  "MappingSecurityRuleId"
// @Param        name query  string  false  "MappingSecurityRuleName"
// @Param        allowedIp query  string  false  "AllowedIpAddress"
// @Param        blockedIp query  string  false  "BlockedIpAddress"
// @Param        createdBeforeAt query  integer  false  "CreatedBeforeAt (Unix timestamp)"
// @Param        createdAfterAt query  integer  false  "CreatedAfterAt (Unix timestamp)"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadMappingSecurityRulesResponse
// @Router       /v1/vhost/mapping/security-rule/ [get]
func (controller *VirtualHostController) ReadMappingSecurityRules(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.ReadMappingSecurityRules(
			requestInputData,
		),
	)
}

// CreateMappingSecurityRule godoc
// @Summary      CreateMappingSecurityRule
// @Description  Create a new mapping security rule.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createMappingSecurityRuleDto body dto.CreateMappingSecurityRule true "Only name is required."
// @Success      201 {object} object{} "MappingSecurityRuleCreated"
// @Router       /v1/vhost/mapping/security-rule/ [post]
func (controller *VirtualHostController) CreateMappingSecurityRule(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if requestInputData["allowedIps"] != nil {
		requestInputData["allowedIps"] = tkPresentation.StringSliceValueObjectParser(
			requestInputData["allowedIps"], tkValueObject.NewCidrBlock,
		)
	}

	if requestInputData["blockedIps"] != nil {
		requestInputData["blockedIps"] = tkPresentation.StringSliceValueObjectParser(
			requestInputData["blockedIps"], tkValueObject.NewCidrBlock,
		)
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.CreateMappingSecurityRule(requestInputData),
	)
}

// UpdateMappingSecurityRule godoc
// @Summary      UpdateMappingSecurityRule
// @Description  Update a mapping security rule.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateMappingSecurityRuleDto body dto.UpdateMappingSecurityRule true "Only id is required."
// @Success      200 {object} object{} "MappingSecurityRuleUpdated"
// @Router       /v1/vhost/mapping/security-rule/ [put]
func (controller *VirtualHostController) UpdateMappingSecurityRule(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if requestInputData["allowedIps"] != nil {
		requestInputData["allowedIps"] = tkPresentation.StringSliceValueObjectParser(
			requestInputData["allowedIps"], tkValueObject.NewCidrBlock,
		)
	}

	if requestInputData["blockedIps"] != nil {
		requestInputData["blockedIps"] = tkPresentation.StringSliceValueObjectParser(
			requestInputData["blockedIps"], tkValueObject.NewCidrBlock,
		)
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.UpdateMappingSecurityRule(requestInputData),
	)
}

// DeleteMappingSecurityRule godoc
// @Summary      DeleteMappingSecurityRule
// @Description  Delete a mapping security rule.
// @Tags         vhosts
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path uint true "MappingSecurityRuleId to delete."
// @Success      200 {object} object{} "MappingSecurityRuleDeleted"
// @Router       /v1/vhost/mapping/security-rule/{id}/ [delete]
func (controller *VirtualHostController) DeleteMappingSecurityRule(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.virtualHostLiaison.DeleteMappingSecurityRule(requestInputData),
	)
}
