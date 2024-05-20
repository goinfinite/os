package vhostInfra

import (
	"errors"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	infraData "github.com/speedianet/os/src/infra/infraData"
)

type VirtualHostCmdRepo struct {
}

func (repo VirtualHostCmdRepo) getAliasConfigFile(
	parentHostname valueObject.Fqdn,
) (valueObject.UnixFilePath, error) {
	vhostFileStr := "/app/conf/nginx/" + parentHostname.String() + ".conf"

	isParentPrimaryVhost := infraHelper.IsPrimaryVirtualHost(
		parentHostname,
	)
	if isParentPrimaryVhost {
		vhostFileStr = "/app/conf/nginx/primary.conf"
	}

	return valueObject.NewUnixFilePath(vhostFileStr)
}

func (repo VirtualHostCmdRepo) createAlias(createDto dto.CreateVirtualHost) error {
	vhostFile, err := repo.getAliasConfigFile(*createDto.ParentHostname)
	if err != nil {
		return errors.New("GetAliasConfigFileFailed")
	}
	vhostFileStr := vhostFile.String()

	hostnameStr := createDto.Hostname.String()

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		`/server_name/ s/;$/ `+hostnameStr+` www.`+hostnameStr+`;/`,
		vhostFileStr,
	)
	if err != nil {
		return errors.New("CreateAliasFailed")
	}

	// TODO: Regenerate cert for primary domain to include new alias

	return infraHelper.ReloadWebServer()
}

func (repo VirtualHostCmdRepo) Create(createDto dto.CreateVirtualHost) error {
	err := infraHelper.ValidateWebServerConfig()
	if err != nil {
		return err
	}

	hostnameStr := createDto.Hostname.String()

	if createDto.Type.String() == "alias" {
		return repo.createAlias(createDto)
	}

	publicDir := infraData.GlobalConfigs.PrimaryPublicDir + "/" + hostnameStr
	certPath := infraData.GlobalConfigs.PkiConfDir + "/" + hostnameStr + ".crt"
	keyPath := infraData.GlobalConfigs.PkiConfDir + "/" + hostnameStr + ".key"
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
	err = infraHelper.UpdateFile(
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

	err = infraHelper.CreateSelfSignedSsl(infraData.GlobalConfigs.PkiConfDir, hostnameStr)
	if err != nil {
		return errors.New("GenerateSelfSignedCertFailed")
	}

	directories := []string{
		publicDir,
		"/app/conf/nginx",
		infraData.GlobalConfigs.PkiConfDir,
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

	return infraHelper.ReloadWebServer()
}

func (repo VirtualHostCmdRepo) deleteAlias(vhost entity.VirtualHost) error {
	vhostFile, err := repo.getAliasConfigFile(*vhost.ParentHostname)
	if err != nil {
		return errors.New("GetAliasConfigFileFailed")
	}
	vhostFileStr := vhostFile.String()

	hostnameStr := vhost.Hostname.String()

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		`/server_name/ s/ `+hostnameStr+` www.`+hostnameStr+`//`,
		vhostFileStr,
	)
	if err != nil {
		return errors.New("DeleteAliasFailed")
	}

	return infraHelper.ReloadWebServer()
}

func (repo VirtualHostCmdRepo) Delete(vhost entity.VirtualHost) error {
	err := infraHelper.ValidateWebServerConfig()
	if err != nil {
		return err
	}

	hostnameStr := vhost.Hostname.String()
	if vhost.Type.String() == "alias" {
		return repo.deleteAlias(vhost)
	}

	_, err = infraHelper.RunCmd(
		"rm",
		"-rf",
		infraData.GlobalConfigs.PrimaryPublicDir+"/"+hostnameStr,
		"/app/conf/nginx/"+hostnameStr+".conf",
		"/app/conf/pki/"+hostnameStr+".crt",
		"/app/conf/pki/"+hostnameStr+".key",
		"/app/conf/nginx/mapping/"+hostnameStr+".conf",
	)
	if err != nil {
		return errors.New("DeleteVirtualHostFailed")
	}

	return infraHelper.ReloadWebServer()
}
