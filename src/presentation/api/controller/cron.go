package apiController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/labstack/echo/v4"
)

type CronController struct {
	cronLiaison *liaison.CronLiaison
}

func NewCronController(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *CronController {
	return &CronController{
		cronLiaison: liaison.NewCronLiaison(trailDbSvc),
	}
}

// ReadCrons	 godoc
// @Summary      ReadCrons
// @Description  List crons.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id query  uint  false  "Id"
// @Param        comment query  string  false  "Comment"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadCronsResponse
// @Router       /v1/cron/ [get]
func (controller *CronController) Read(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(c, controller.cronLiaison.Read(requestInputData))
}

// CreateCron    godoc
// @Summary      CreateCron
// @Description  Create a new cron.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createCronDto 	  body    dto.CreateCron  true  "comment is optional."
// @Success      201 {object} object{} "CronCreated"
// @Router       /v1/cron/ [post]
func (controller *CronController) Create(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.cronLiaison.Create(requestInputData),
	)
}

// UpdateCron godoc
// @Summary      UpdateCron
// @Description  Update a cron.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateCronDto 	  body dto.UpdateCron  true  "Only id is required."
// @Success      200 {object} object{} "CronUpdated message"
// @Router       /v1/cron/ [put]
func (controller *CronController) Update(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.cronLiaison.Update(requestInputData),
	)
}

// DeleteCron	 godoc
// @Summary      DeleteCron
// @Description  Delete a cron.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        cronId 	  path   string  true  "CronId to delete."
// @Success      200 {object} object{} "CronDeleted"
// @Router       /v1/cron/{cronId}/ [delete]
func (controller *CronController) Delete(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.cronLiaison.Delete(requestInputData),
	)
}
