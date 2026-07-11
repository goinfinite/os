package cliMiddleware

import (
	"fmt"
	"os"
	"os/user"
)

func PreventRootless() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("ReadCurrentUserError: ", err)
		os.Exit(1)
	}

	if currentUser.Username != "root" {
		fmt.Println("InfiniteOsMustBeRunAsRoot")
		os.Exit(1)
	}
}
