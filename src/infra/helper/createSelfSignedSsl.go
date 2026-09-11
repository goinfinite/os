package infraHelper

import (
	"errors"
	"log/slog"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

func CreateSelfSignedSsl(
	dirPath tkValueObject.UnixAbsoluteFilePath,
	vhostHostname tkValueObject.Fqdn,
	aliasesHostname []tkValueObject.Fqdn,
) error {
	dirPathStr := dirPath.String()
	vhostHostnameStr := vhostHostname.String()

	rawCertFilePath := dirPathStr + "/" + vhostHostnameStr + ".crt"
	certFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(rawCertFilePath, false)
	if err != nil {
		return errors.New("InvalidCertFilePath: " + err.Error())
	}
	certFilePathStr := certFilePath.String()

	rawKeyFilePath := dirPathStr + "/" + vhostHostnameStr + ".key"
	keyFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(rawKeyFilePath, false)
	if err != nil {
		return errors.New("InvalidKeyFilePath: " + err.Error())
	}
	keyFilePathStr := keyFilePath.String()

	fileClerk := tkInfra.FileClerk{}
	if fileClerk.FileExists(certFilePathStr) && fileClerk.FileExists(keyFilePathStr) {
		slog.Debug(
			"SkippingSelfSignedSslAlreadyExists",
			slog.String("certFilePath", certFilePathStr),
			slog.String("keyFilePath", keyFilePathStr),
		)
		return nil
	}

	synth := &tkInfra.Synthesizer{}
	certPem, keyPem, err := synth.SelfSignedCertificatePairPemFactory(
		&vhostHostname, aliasesHostname,
	)
	if err != nil {
		return errors.New("SelfSignedCertificateGenerationError: " + err.Error())
	}

	webServerUsername := tkValueObject.UnixUsername(infraEnvs.PhpWebServerUsername)
	symlinkPolicy := tkInfra.FileClerkSymlinkPolicyResolve
	overwritePolicy := tkInfra.FileClerkOverwritePolicyReplace

	privateKeyPermissions := os.FileMode(0600)
	err = fileClerk.UpsertFile(tkInfra.FileUpsertSettings{
		FilePath:                keyFilePath,
		Permissions:             &privateKeyPermissions,
		SymlinkPolicy:           &symlinkPolicy,
		OverwritePolicy:         &overwritePolicy,
		OwnerUsername:           &webServerUsername,
		TrustedDirOwnerUsername: &webServerUsername,
	}, []byte(keyPem))
	if err != nil {
		return errors.New("WritePrivateKeyError: " + err.Error())
	}

	certificatePermissions := os.FileMode(0644)
	err = fileClerk.UpsertFile(tkInfra.FileUpsertSettings{
		FilePath:                certFilePath,
		Permissions:             &certificatePermissions,
		SymlinkPolicy:           &symlinkPolicy,
		OverwritePolicy:         &overwritePolicy,
		OwnerUsername:           &webServerUsername,
		TrustedDirOwnerUsername: &webServerUsername,
	}, []byte(certPem))
	if err != nil {
		return errors.New("WriteCertificateError: " + err.Error())
	}

	return nil
}
