package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	apiHelper "github.com/speedianet/sam/src/presentation/api/helper"
)

// GetCrons	 godoc
// @Summary      GetCrons
// @Description  List Crons.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Cron
// @Router       /cron/ [get]
func GetCronsController(c echo.Context) error {
	cronQueryRepo := infra.CronQueryRepo{}
	cronsList, err := useCase.GetCrons(cronQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, cronsList)
}

// AddCron    godoc
// @Summary      AddNewCron
// @Description  Add a new cron.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addCronDto 	  body    dto.AddCron  true  "NewCron"
// @Success      201 {object} object{} "CronCreated"
// @Router       /cron/ [post]
func AddCronController(c echo.Context) error {
	requiredParams := []string{"schedule", "command"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	var cronCommentPtr *valueObject.CronComment
	if requestBody["comment"] != nil {
		cronComment := valueObject.NewCronCommentPanic(requestBody["comment"].(string))
		cronCommentPtr = &cronComment
	}

	addCronDto := dto.NewAddCron(
		valueObject.NewCronSchedulePanic(requestBody["schedule"].(string)),
		valueObject.NewUnixCommandPanic(requestBody["command"].(string)),
		cronCommentPtr,
	)

	cronCmdRepo, err := infra.NewCronCmdRepo()
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	err = useCase.AddCron(
		cronCmdRepo,
		addCronDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "CronCreated")
}

// UpdateCron godoc
// @Summary      UpdateCron
// @Description  Update an cron.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateCronDto 	  body dto.UpdateCron  true  "UpdateCron"
// @Success      200 {object} object{} "CronUpdated message"
// @Router       /cron/ [put]
func UpdateCronController(c echo.Context) error {
	requiredParams := []string{"id"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	var cronSchedulePtr *valueObject.CronSchedule
	if requestBody["schedule"] != nil {
		cronSchedule := valueObject.NewCronSchedulePanic(requestBody["schedule"].(string))
		cronSchedulePtr = &cronSchedule
	}

	var cronCommandPtr *valueObject.UnixCommand
	if requestBody["command"] != nil {
		cronCommand := valueObject.NewUnixCommandPanic(requestBody["command"].(string))
		cronCommandPtr = &cronCommand
	}

	var cronCommentPtr *valueObject.CronComment
	if requestBody["comment"] != nil {
		cronComment := valueObject.NewCronCommentPanic(requestBody["comment"].(string))
		cronCommentPtr = &cronComment
	}

	updateCronDto := dto.NewUpdateCron(
		valueObject.NewCronIdPanic(requestBody["id"]),
		cronSchedulePtr,
		cronCommandPtr,
		cronCommentPtr,
	)

	cronQueryRepo := infra.CronQueryRepo{}
	cronCmdRepo, err := infra.NewCronCmdRepo()
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	err = useCase.UpdateCron(
		cronQueryRepo,
		cronCmdRepo,
		updateCronDto,
	)

	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "CronUpdated")
}

// DeleteCron	 godoc
// @Summary      DeleteCron
// @Description  Delete an cron.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        cronId 	  path   string  true  "CronId"
// @Success      200 {object} object{} "CronDeleted"
// @Router       /cron/{cronId}/ [delete]
func DeleteCronController(c echo.Context) error {
	cronId := valueObject.NewCronIdPanic(c.Param("cronId"))

	cronQueryRepo := infra.CronQueryRepo{}
	cronCmdRepo, err := infra.NewCronCmdRepo()
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	err = useCase.DeleteCron(
		cronQueryRepo,
		cronCmdRepo,
		cronId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "CronDeleted")
}
