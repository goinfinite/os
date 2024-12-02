package apiController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type CronController struct {
	cronService *service.CronService
}

func NewCronController(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) (*CronController, error) {
	cronService, err := service.NewCronService(trailDbSvc)
	if err != nil {
		return nil, err
	}

	return &CronController{
		cronService: cronService,
	}, nil
}

// ReadCrons	 godoc
// @Summary      ReadCrons
// @Description  List crons.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Cron
// @Router       /v1/cron/ [get]
func (controller *CronController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.cronService.Read())
}

// CreateCron    godoc
// @Summary      CreateNewCron
// @Description  Create a new cron.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createCronDto 	  body    dto.CreateCron  true  "comment is optional."
// @Success      201 {object} object{} "CronCreated"
// @Router       /v1/cron/ [post]
func (controller *CronController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.cronService.Create(requestBody),
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
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.cronService.Update(requestBody),
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
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.cronService.Delete(requestBody),
	)
}
