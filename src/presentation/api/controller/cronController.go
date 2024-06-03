package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	cronInfra "github.com/speedianet/os/src/infra/cron"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

// GetCrons	 godoc
// @Summary      GetCrons
// @Description  List crons.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Cron
// @Router       /v1/cron/ [get]
func GetCronsController(c echo.Context) error {
	cronQueryRepo := cronInfra.CronQueryRepo{}
	cronsList, err := useCase.GetCrons(cronQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, cronsList)
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
func CreateCronController(c echo.Context) error {
	requiredParams := []string{"schedule", "command"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	var cronCommentPtr *valueObject.CronComment
	if requestBody["comment"] != nil {
		rawComment, assertOk := requestBody["comment"].(string)
		if !assertOk {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "InvalidComment")
		}

		if len(rawComment) > 0 {
			cronComment := valueObject.NewCronCommentPanic(rawComment)
			cronCommentPtr = &cronComment
		}
	}

	createCronDto := dto.NewCreateCron(
		valueObject.NewCronSchedulePanic(requestBody["schedule"].(string)),
		valueObject.NewUnixCommandPanic(requestBody["command"].(string)),
		cronCommentPtr,
	)

	cronCmdRepo, err := cronInfra.NewCronCmdRepo()
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	err = useCase.CreateCron(
		cronCmdRepo,
		createCronDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "CronCreated")
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

	cronQueryRepo := cronInfra.CronQueryRepo{}
	cronCmdRepo, err := cronInfra.NewCronCmdRepo()
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	err = useCase.UpdateCron(
		cronQueryRepo,
		cronCmdRepo,
		updateCronDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "CronUpdated")
}

// DeleteCron	 godoc
// @Summary      DeleteCron
// @Description  Delete a cron.
// @Tags         cron
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        cronId 	  path   string  true  "Cron ID that will be deleted."
// @Success      200 {object} object{} "CronDeleted"
// @Router       /v1/cron/{cronId}/ [delete]
func DeleteCronController(c echo.Context) error {
	cronId := valueObject.NewCronIdPanic(c.Param("cronId"))

	cronQueryRepo := cronInfra.CronQueryRepo{}
	cronCmdRepo, err := cronInfra.NewCronCmdRepo()
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	err = useCase.DeleteCron(
		cronQueryRepo,
		cronCmdRepo,
		cronId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "CronDeleted")
}
