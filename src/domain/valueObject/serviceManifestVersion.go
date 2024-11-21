package valueObject

import (
	"errors"
	"slices"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

type ServiceManifestVersion string

var validServiceManifestVersions = []string{
	"v1",
}

func NewServiceManifestVersion(value interface{}) (
	version ServiceManifestVersion, err error,
) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return version, errors.New("ServiceManifestVersionMustBeString")
	}
	stringValue = strings.ToLower(stringValue)

	if !slices.Contains(validServiceManifestVersions, stringValue) {
		return version, errors.New("InvalidServiceManifestVersion")
	}

	return ServiceManifestVersion(stringValue), nil
}

func (vo ServiceManifestVersion) String() string {
	return string(vo)
}
