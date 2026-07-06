package vhostInfra

import (
	"errors"
	"log/slog"
	"os"
	"time"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

type VirtualHostHelpers struct{}

func NewVirtualHostHelpers() *VirtualHostHelpers {
	return &VirtualHostHelpers{}
}

func (helpers *VirtualHostHelpers) parsePrimaryConfHostname() (
	primaryHostname tkValueObject.Fqdn, err error,
) {
	rawServerNameHostname, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-n",
			`/^\s\{0,20\}server_name\s/s/^\s\{0,20\}server_name\s\{1,255\}\([^; ]\{1,255\}\).\{0,1024\}$/\1/p`,
			infraEnvs.PrimaryVirtualHostConfPath,
		},
	}).Run()
	if err != nil {
		return primaryHostname, err
	}

	if rawServerNameHostname == "" {
		return primaryHostname, errors.New("PrimaryServerNameNotFound")
	}

	return tkValueObject.NewFqdn(rawServerNameHostname)
}

func (helpers *VirtualHostHelpers) ReadPrimaryVirtualHostHostname() (
	primaryHostname tkValueObject.Fqdn, err error,
) {
	primaryHostFromEnv := os.Getenv("PRIMARY_VHOST")
	if primaryHostFromEnv != "" {
		return tkValueObject.NewFqdn(primaryHostFromEnv)
	}

	hostnameFromConf, parseErr := helpers.parsePrimaryConfHostname()
	if parseErr == nil {
		return hostnameFromConf, nil
	}
	slog.Debug(
		"ParsePrimaryConfHostnameFail",
		slog.String("error", parseErr.Error()),
	)

	primaryHostFromShell, shellErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "hostname",
		Args:    []string{"-f"},
	}).Run()
	if shellErr != nil {
		return primaryHostname, shellErr
	}

	return tkValueObject.NewFqdn(primaryHostFromShell)
}

func (helpers *VirtualHostHelpers) IsPrimaryVirtualHost(vhost tkValueObject.Fqdn) bool {
	primaryVhost, err := helpers.ReadPrimaryVirtualHostHostname()
	if err != nil {
		slog.Error(
			"ReadPrimaryVirtualHostHostnameError",
			slog.String("error", err.Error()),
		)
		return false
	}

	return vhost == primaryVhost
}

func (helpers *VirtualHostHelpers) ValidateWebServerConfig() error {
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:           infraEnvs.WebServerBinaryPath + " -t",
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return errors.New("WebServerConfigValidationError: " + err.Error())
	}

	return nil
}

func (helpers *VirtualHostHelpers) ReloadWebServer() error {
	err := helpers.ValidateWebServerConfig()
	if err != nil {
		return errors.New("WebServerConfigTestFail: " + err.Error())
	}

	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command: infraEnvs.WebServerBinaryPath + " -s reload -c " +
			infraEnvs.WebServerMainConfPath,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return errors.New("WebServerReloadFail: " + err.Error())
	}

	time.Sleep(1 * time.Second)

	return nil
}
