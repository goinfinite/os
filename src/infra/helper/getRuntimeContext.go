package infraHelper

import (
	"os/exec"
	"strings"

	"github.com/speedianet/sam/src/domain/valueObject"
)

func isVm() bool {
	out, err := exec.Command("ps", "-p", "1", "-o", "comm=").Output()
	if err != nil {
		return false
	}

	output := strings.TrimSpace(string(out))
	return output == "systemd"
}

func GetRuntimeContext() (valueObject.RuntimeContext, error) {
	if isVm() {
		return valueObject.NewRuntimeContext("vm")
	}
	return valueObject.NewRuntimeContext("container")
}
