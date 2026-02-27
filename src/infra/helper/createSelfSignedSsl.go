package infraHelper

import (
	"errors"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

func CreateSelfSignedSsl(
	dirPath tkValueObject.UnixAbsoluteFilePath,
	vhostHostname tkValueObject.Fqdn,
	aliasesHostname []tkValueObject.Fqdn,
) error {
	synth := &tkInfra.Synthesizer{}
	certPem, keyPem, err := synth.SelfSignedCertificatePairPemFactory(
		&vhostHostname, aliasesHostname,
	)
	if err != nil {
		return errors.New("SelfSignedCertificateGenerationError: " + err.Error())
	}

	vhostHostnameStr := vhostHostname.String()
	dirPathStr := dirPath.String()

	shouldOverwrite := true
	sslFileClerk := tkInfra.FileClerk{}
	err = sslFileClerk.UpdateFileContent(
		dirPathStr+"/"+vhostHostnameStr+".key",
		keyPem, shouldOverwrite,
	)
	if err != nil {
		return errors.New("WritePrivateKeyError: " + err.Error())
	}

	err = sslFileClerk.UpdateFileContent(
		dirPathStr+"/"+vhostHostnameStr+".crt",
		certPem, shouldOverwrite,
	)
	if err != nil {
		return errors.New("WriteCertificateError: " + err.Error())
	}

	return nil
}
