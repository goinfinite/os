package infraHelper

func GetOsRelease() (string, error) {
	out, err := RunCmd("lsb_release", "-cs")
	if err != nil {
		return "", err
	}
	return out, nil
}
