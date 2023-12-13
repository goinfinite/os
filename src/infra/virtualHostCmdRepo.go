package infra

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type VirtualHostCmdRepo struct {
}

func (repo VirtualHostCmdRepo) Add(addDto dto.AddVirtualHost) error {
	hostnameStr := addDto.Hostname.String()

	publicDir := "/app/html/" + hostnameStr
	certPath := "/app/conf/pki/" + hostnameStr + ".crt"
	keyPath := "/app/conf/pki/" + hostnameStr + ".key"
	mappingFilePath := "/app/conf/nginx/mapping/" + hostnameStr + ".conf"

	nginxConf := `server {
    listen 80;
    listen 443 ssl;
    server_name ` + hostnameStr + ` www.` + hostnameStr + `;

    root ` + publicDir + `;

    ssl_certificate ` + certPath + `;
    ssl_certificate_key ` + keyPath + `;

    access_log /app/logs/nginx/` + hostnameStr + `_access.log combined buffer=512k flush=1m;
    error_log /app/logs/nginx/` + hostnameStr + `_error.log warn;

    include /etc/nginx/std.conf;
    include ` + mappingFilePath + `;
}
`
	err := infraHelper.UpdateFile(
		"/app/conf/nginx/"+hostnameStr+".conf",
		nginxConf,
		true,
	)
	if err != nil {
		return errors.New("CreateNginxConfFileFailed")
	}

	err = infraHelper.UpdateFile(
		mappingFilePath,
		"",
		true,
	)
	if err != nil {
		return errors.New("CreateMappingFileFailed")
	}

	err = infraHelper.MakeDir(publicDir)
	if err != nil {
		return errors.New("MakePublicHtmlDirFailed")
	}

	_, err = infraHelper.RunCmd(
		"openssl",
		"req",
		"-x509",
		"-nodes",
		"-days",
		"365",
		"-newkey",
		"rsa:2048",
		"-keyout",
		keyPath,
		"-out",
		certPath,
		"-subj",
		"/C=US/ST=California/L=LosAngeles/O=Acme/CN="+hostnameStr,
	)
	if err != nil {
		return errors.New("GenerateSelfSignedCertFailed")
	}

	directories := []string{
		publicDir,
		"/app/conf/nginx",
		"/app/conf/pki",
	}
	for _, directory := range directories {
		_, err = infraHelper.RunCmd(
			"chown",
			"-R",
			"nobody:nogroup",
			directory,
		)
		if err != nil {
			return errors.New("ChownNecessaryDirectoriesFailed")
		}
	}

	_, err = infraHelper.RunCmd(
		"nginx",
		"-t",
	)
	if err != nil {
		return errors.New("NginxConfigTestFailed")
	}

	_, err = infraHelper.RunCmd(
		"nginx",
		"-s",
		"reload",
	)
	if err != nil {
		return errors.New("NginxReloadFailed")
	}

	return nil
}
