package voHelper

func InterfaceToUint16(input interface{}) (uint16Value uint16, err error) {
	uintValue, err := InterfaceToUint(input)
	if err != nil {
		return uint16Value, err
	}

	return uint16(uintValue), err
}
