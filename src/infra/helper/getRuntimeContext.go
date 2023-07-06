package infraHelper

import (
	"bufio"
	"os"
	"strings"

	"github.com/speedianet/sam/src/domain/valueObject"
)

func isContainer() bool {
	file, err := os.Open("/proc/self/cgroup")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "docker") || strings.Contains(line, "podman") {
			return true
		}
	}

	return false
}

func GetRuntimeContext() (valueObject.RuntimeContext, error) {
	if isContainer() {
		return valueObject.NewRuntimeContext("container")
	}
	return valueObject.NewRuntimeContext("vm")
}
