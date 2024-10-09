package cliMiddleware

import (
	"fmt"
	"os"
	"os/user"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

func PreventRootless() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("ReadCurrentUserError: ", err)
		os.Exit(1)
	}

	if isDevMode, _ := voHelper.InterfaceToBool(os.Getenv("DEV_MODE")); isDevMode {
		fmt.Println("BinaryCompiledSuccessfully")
		os.Exit(0)
	}

	if currentUser.Username != "root" {
		fmt.Println("InfiniteOsMustBeRunAsRoot")
		os.Exit(1)
	}
}
