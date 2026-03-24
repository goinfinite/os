package cliMiddleware

import (
	"fmt"
	"os"
	"os/user"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

func PreventRootless() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("ReadCurrentUserError: ", err)
		os.Exit(1)
	}

	if isDevMode, _ := tkVoUtil.InterfaceToBool(os.Getenv("DEV_MODE")); isDevMode {
		fmt.Println("BinaryCompiledSuccessfully")
		os.Exit(0)
	}

	if currentUser.Username != "root" {
		fmt.Println("InfiniteOsMustBeRunAsRoot")
		os.Exit(1)
	}
}
