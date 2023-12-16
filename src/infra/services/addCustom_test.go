package servicesInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestAddCustom(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("AddCustomService", func(t *testing.T) {
		t.Skip("SkipAddCustomServiceTest")

		portBinding, err := valueObject.NewPortBindingFromString(
			"8000/http",
		)
		if err != nil {
			t.Errorf("NewPortBindingFromStringFailed : %v", err)
			return
		}

		dto := dto.NewAddCustomService(
			valueObject.NewServiceNamePanic("python-ws"),
			valueObject.NewServiceTypePanic("webserver"),
			valueObject.NewUnixCommandPanic("python3 -m http.server"),
			nil,
			[]valueObject.PortBinding{portBinding},
		)

		err = AddCustom(dto)
		if err != nil {
			t.Errorf("AddCustomServiceFailed : %v", err)
			return
		}
	})
}
