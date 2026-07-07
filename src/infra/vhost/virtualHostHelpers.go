package vhostInfra

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

type VirtualHostHelpers struct{}

func NewVirtualHostHelpers() *VirtualHostHelpers {
	return &VirtualHostHelpers{}
}

func (helpers *VirtualHostHelpers) ReadPrimaryVirtualHostHostnameFromWebServerConf() (
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
	primaryHostFromEnv := os.Getenv(infraEnvs.PrimaryVirtualHostEnvKey)
	if primaryHostFromEnv != "" {
		return tkValueObject.NewFqdn(primaryHostFromEnv)
	}
	slog.Debug("PrimaryVirtualHostEnvValueNotFound")

	hostnameFromConf, parseErr := helpers.ReadPrimaryVirtualHostHostnameFromWebServerConf()
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

func (helpers *VirtualHostHelpers) UpdateWebServerWorkerCount(
	cpuCoresStr string,
	servicesCmdRepo repository.ServicesCmdRepo,
) error {
	_, sedErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-i", "-e",
			"s/^worker_processes.*/worker_processes " + cpuCoresStr + ";/g",
			infraEnvs.WebServerMainConfPath,
		},
	}).Run()
	if sedErr != nil {
		return errors.New("UpdateNginxWorkersCountFailed: " + sedErr.Error())
	}

	serviceName, _ := valueObject.NewServiceName("nginx")
	restartErr := servicesCmdRepo.Restart(serviceName)
	if restartErr != nil {
		return errors.New("RestartNginxFailed: " + restartErr.Error())
	}

	return nil
}

func (helpers *VirtualHostHelpers) UpdateWebServerPrimaryVirtualHost(
	newHostname tkValueObject.Fqdn,
) error {
	currentHostname, readErr := helpers.ReadPrimaryVirtualHostHostnameFromWebServerConf()
	if readErr != nil {
		return errors.New(
			"ReadPrimaryVirtualHostHostnameFromWebServerConfFailed: " + readErr.Error(),
		)
	}

	if currentHostname == newHostname {
		return nil
	}

	currentHostnameStr := currentHostname.String()
	grepCurrentHostname := strings.ReplaceAll(currentHostnameStr, ".", "\\.")
	rawGrepOutput, grepErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "grep",
		Args: []string{
			"-n", "--",
			"server_name.*" + grepCurrentHostname,
			infraEnvs.PrimaryVirtualHostConfPath,
		},
	}).Run()
	if grepErr != nil {
		return errors.New("FindHostnameLineFailed: " + grepErr.Error())
	}

	if rawGrepOutput == "" {
		return errors.New("HostnameLineNotFound")
	}

	rawGrepOutputParts := strings.SplitN(rawGrepOutput, ":", 2)
	if len(rawGrepOutputParts) < 2 {
		return errors.New("InvalidGrepOutputFormat")
	}

	rawHostnameLineNumStr := rawGrepOutputParts[0]

	hostnameLineNum, parseErr := strconv.Atoi(rawHostnameLineNumStr)
	if parseErr != nil {
		return errors.New("ParseHostnameLineNumFailed: " + parseErr.Error())
	}

	sedCurrentHostname := strings.ReplaceAll(currentHostnameStr, ".", "\\.")
	sedNewHostname := strings.ReplaceAll(newHostname.String(), ".", "\\.")
	_, updateErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-i",
			strconv.Itoa(hostnameLineNum) + "s/" + sedCurrentHostname + "/" +
				sedNewHostname + "/g",
			infraEnvs.PrimaryVirtualHostConfPath,
		},
	}).Run()
	if updateErr != nil {
		return errors.New(
			"UpdatePrimaryVhostServerNameFailed: " + updateErr.Error(),
		)
	}

	return nil
}
