package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestCreateCustom(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("CreateCustomService", func(t *testing.T) {
		t.Skip("SkipCreateCustomServiceTest")

		portBinding, err := valueObject.NewPortBindingFromString(
			"8000/http",
		)
		if err != nil {
			t.Errorf("NewPortBindingFromStringFailed : %v", err)
			return
		}

		dto := dto.NewCreateCustomService(
			valueObject.NewServiceNamePanic("python-ws"),
			valueObject.NewServiceTypePanic("webserver"),
			valueObject.NewUnixCommandPanic("python3 -m http.server"),
			nil,
			[]valueObject.PortBinding{portBinding},
			true,
		)

		err = CreateCustom(dto)
		if err != nil {
			t.Errorf("CreateCustomServiceFailed : %v", err)
			return
		}
	})
}
