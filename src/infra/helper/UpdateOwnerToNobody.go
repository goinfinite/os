package infraHelper

func UpdateOwnerToNobody(
	inodePath string, recursively bool, includeSymlinks bool,
) error {
	chownCmd := "chown "
	if recursively {
		chownCmd += "-R "
	}

	chownCmd += "nobody:nogroup "
	if includeSymlinks {
		chownCmd += "-L "
	}

	chownCmd += inodePath
	_, err := RunCmdWithSubShell(chownCmd)
	if err != nil {
		return err
	}

	return nil
}
