package voHelper

func InterfaceToUint16(input interface{}) (output uint16, err error) {
	uintValue, err := InterfaceToUint(input)
	if err != nil {
		return output, err
	}

	return uint16(uintValue), err
}
