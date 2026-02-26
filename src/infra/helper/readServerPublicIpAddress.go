package infraHelper

import (
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

func ReadServerPublicIpAddress() (tkValueObject.IpAddress, error) {
	return tkInfra.ReadServerPublicIpAddress()
}
