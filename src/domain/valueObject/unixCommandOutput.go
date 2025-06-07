package valueObject

import (
	"errors"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

type UnixCommandOutput string

func NewUnixCommandOutput(rawValue any) (cmdOutput UnixCommandOutput, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(rawValue)
	if err != nil {
		return cmdOutput, errors.New("UnixCommandOutputMustBeString")
	}

	if len(stringValue) > 4096 {
		stringValue = stringValue[:4096]
	}

	return UnixCommandOutput(stringValue), nil
}

func (vo UnixCommandOutput) String() string {
	return string(vo)
}
