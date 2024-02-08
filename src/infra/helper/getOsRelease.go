package infraHelper

func GetOsRelease() (string, error) {
	return RunCmd("lsb_release", "-cs")
}
