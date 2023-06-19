package infra

import (
	"testing"
	"time"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/valueObject"
)

func TestAuthCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()

	t.Run("GetSessionToken", func(t *testing.T) {
		authCmdRepo := AuthCmdRepo{}
		token := authCmdRepo.GenerateSessionToken(
			valueObject.UserId(1000),
			valueObject.UnixTime(
				time.Now().Add(3*time.Hour).Unix(),
			),
			valueObject.NewIpAddressPanic("127.0.0.1"),
		)
		if token.TokenStr == "" {
			t.Errorf("Expected token not to be empty")
		}
	})
}
