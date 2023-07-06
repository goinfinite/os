package valueObject

import "errors"

type RuntimeContext string

const (
	container RuntimeContext = "container"
	vm        RuntimeContext = "vm"
	bareMetal RuntimeContext = "bareMetal"
)

func NewRuntimeContext(value string) (RuntimeContext, error) {
	rc := RuntimeContext(value)
	if !rc.isValid() {
		return "", errors.New("InvalidRuntimeContext")
	}
	return rc, nil
}

func NewRuntimeContextPanic(value string) RuntimeContext {
	rc := RuntimeContext(value)
	if !rc.isValid() {
		panic("InvalidRuntimeContext")
	}
	return rc
}

func (rc RuntimeContext) isValid() bool {
	switch rc {
	case container, vm, bareMetal:
		return true
	default:
		return false
	}
}

func (rc RuntimeContext) IsContainer() bool {
	return rc == container
}

func (rc RuntimeContext) String() string {
	return string(rc)
}
