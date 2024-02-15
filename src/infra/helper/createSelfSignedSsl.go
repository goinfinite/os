package infraHelper

import (
	"fmt"
)

func CreateSelfSignedSsl(virtualHost string) error {
	vhostCertKeyFilePath := "/app/conf/pki/" + virtualHost + ".key"
	vhostCertFilePath := "/app/conf/pki/" + virtualHost + ".crt"

	_, err := RunCmd(
		"openssl",
		"req",
		"-x509",
		"-nodes",
		"-days",
		"365",
		"-newkey",
		"rsa:2048",
		"-keyout",
		vhostCertKeyFilePath,
		"-out",
		vhostCertFilePath,
		"-subj",
		"/C=US/ST=California/L=LosAngeles/O=Acme/CN="+virtualHost,
	)
	if err != nil {
		return fmt.Errorf("ReplaceWithSelfSignedFailed (%s): %s", virtualHost, err.Error())
	}

	return nil
}
