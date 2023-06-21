package restApiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	restApiHelper "github.com/speedianet/sam/src/presentation/api/helper"
)

// AuthLogin godoc
// @Summary      AddNewUser
// @Description  Add a new user.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addUserDto 	  body    dto.AddUser  true  "New user details"
// @Success      201 {object} object{} "UserCreated"
// @Router       /user/ [post]
func AddUserController(c echo.Context) error {
	requiredParams := []string{"username", "password"}
	requestBody, _ := restApiHelper.GetRequestBody(c)

	restApiHelper.CheckMissingParams(requestBody, requiredParams)

	addUserDto := dto.NewAddUser(
		valueObject.NewUsernamePanic(requestBody["username"].(string)),
		valueObject.NewPasswordPanic(requestBody["password"].(string)),
	)

	accQueryRepo := infra.AccQueryRepo{}
	accCmdRepo := infra.AccCmdRepo{}

	useCase.AddUser(
		accQueryRepo,
		accCmdRepo,
		addUserDto,
	)

	return restApiHelper.ResponseWrapper(c, http.StatusCreated, "UserCreated")
}

// AuthLogin godoc
// @Summary      DeleteUser
// @Description  Delete an user.
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        userId 	  path   string  true  "UserId"
// @Success      200 {object} object{} "UserDeleted"
// @Router       /user/{userId}/ [delete]
func DeleteUserController(c echo.Context) error {
	userId := valueObject.NewUserIdFromStringPanic(c.Param("userId"))

	accQueryRepo := infra.AccQueryRepo{}
	accCmdRepo := infra.AccCmdRepo{}

	useCase.DeleteUser(
		accQueryRepo,
		accCmdRepo,
		userId,
	)

	return restApiHelper.ResponseWrapper(c, http.StatusOK, "UserDeleted")
}
