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
	infraHelper "github.com/goinfinite/os/src/infra/helper"
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
		slog.String("err", parseErr.Error()),
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
			slog.String("err", err.Error()),
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

func (helpers *VirtualHostHelpers) findWebServerConfLineNumber(
	confPath, searchPattern string,
) (lineNumber int, err error) {
	rawGrepOutput, grepErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "grep",
		Args: []string{
			"-n", "-E", "--", searchPattern, confPath,
		},
	}).Run()
	if grepErr != nil {
		return lineNumber, errors.New("GrepWebServerConfLineFailed: " + grepErr.Error())
	}

	if rawGrepOutput == "" {
		return lineNumber, errors.New("WebServerConfLineNotFound")
	}

	rawParts := strings.SplitN(rawGrepOutput, ":", 2)
	if len(rawParts) < 2 {
		return lineNumber, errors.New("InvalidGrepOutputFormat")
	}

	lineNumber, parseErr := strconv.Atoi(rawParts[0])
	if parseErr != nil {
		return lineNumber, errors.New("ParseLineNumberFailed: " + parseErr.Error())
	}

	return lineNumber, nil
}

func (helpers *VirtualHostHelpers) replaceWebServerConfLineContent(
	confPath string,
	lineNumber int,
	previousValue, newValue string,
) error {
	_, sedErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-i", "-E",
			strconv.Itoa(lineNumber) + "s|" + previousValue + "|" + newValue + "|",
			confPath,
		},
	}).Run()
	if sedErr != nil {
		return errors.New("ReplaceWebServerConfLineContentFailed: " + sedErr.Error())
	}

	return nil
}

func (helpers *VirtualHostHelpers) replaceVirtualHostServerName(
	currentHostname tkValueObject.Fqdn,
	newHostname tkValueObject.Fqdn,
) error {
	escapedHostname := strings.ReplaceAll(currentHostname.String(), ".", "\\.")

	lineNumber, grepErr := helpers.findWebServerConfLineNumber(
		infraEnvs.PrimaryVirtualHostConfPath,
		// 	server_name app.net;
		"^[[:space:]]*server_name[[:space:]]+"+escapedHostname+"[[:space:];]",
	)
	if grepErr != nil {
		return errors.New("FindServerNameLineFailed: " + grepErr.Error())
	}

	newHostnameStr := newHostname.String()
	sedErr := helpers.replaceWebServerConfLineContent(
		infraEnvs.PrimaryVirtualHostConfPath, lineNumber, escapedHostname,
		strings.ReplaceAll(newHostnameStr, ".", "\\."),
	)
	if sedErr != nil {
		return errors.New("UpdateServerNameFailed: " + sedErr.Error())
	}

	return nil
}

func (helpers *VirtualHostHelpers) replaceVirtualHostSslCertificate(
	newHostname tkValueObject.Fqdn,
) error {
	newHostnameStr := newHostname.String()
	rawCertPath := infraEnvs.PkiConfDir + "/" + newHostnameStr + ".crt"
	newCertPath, certPathErr := tkValueObject.NewUnixAbsoluteFilePath(rawCertPath, false)
	if certPathErr != nil {
		return errors.New("InvalidSslCertPath: " + certPathErr.Error())
	}

	rawKeyPath := infraEnvs.PkiConfDir + "/" + newHostnameStr + ".key"
	newKeyPath, keyPathErr := tkValueObject.NewUnixAbsoluteFilePath(rawKeyPath, false)
	if keyPathErr != nil {
		return errors.New("InvalidSslKeyPath: " + keyPathErr.Error())
	}

	confPath := infraEnvs.PrimaryVirtualHostConfPath

	certLineNum, certGrepErr := helpers.findWebServerConfLineNumber(
		confPath,
		// ssl_certificate /app/conf/pki/app.net.crt;
		"^[[:space:]]*ssl_certificate[[:space:]]",
	)
	if certGrepErr != nil {
		return errors.New("FindSslCertLineFailed: " + certGrepErr.Error())
	}

	keyLineNum, keyGrepErr := helpers.findWebServerConfLineNumber(
		confPath,
		// ssl_certificate_key /app/conf/pki/app.net.key;
		"^[[:space:]]*ssl_certificate_key[[:space:]]",
	)
	if keyGrepErr != nil {
		return errors.New("FindSslKeyLineFailed: " + keyGrepErr.Error())
	}

	certSedErr := helpers.replaceWebServerConfLineContent(
		confPath, certLineNum,
		// ssl_certificate /old/path.crt; → ssl_certificate /new/path.crt;
		"ssl_certificate[[:space:]]+[^;]+",
		"ssl_certificate "+newCertPath.String()+";",
	)
	if certSedErr != nil {
		return errors.New("UpdateSslCertPathFailed: " + certSedErr.Error())
	}

	keySedErr := helpers.replaceWebServerConfLineContent(
		confPath, keyLineNum,
		// ssl_certificate_key /old/path.key; → ssl_certificate_key /new/path.key;
		"ssl_certificate_key[[:space:]]+[^;]+",
		"ssl_certificate_key "+newKeyPath.String()+";",
	)
	if keySedErr != nil {
		return errors.New("UpdateSslKeyPathFailed: " + keySedErr.Error())
	}

	return nil
}

func (helpers *VirtualHostHelpers) UpdateWebServerPrimaryVirtualHost(
	newHostname tkValueObject.Fqdn,
	aliasesHostnames []tkValueObject.Fqdn,
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

	serverNameErr := helpers.replaceVirtualHostServerName(
		currentHostname, newHostname,
	)
	if serverNameErr != nil {
		return serverNameErr
	}

	pkiConfDir, pkiErr := tkValueObject.NewUnixAbsoluteFilePath(infraEnvs.PkiConfDir, false)
	if pkiErr != nil {
		return errors.New("PkiConfDirNotFound: " + pkiErr.Error())
	}

	createCertErr := infraHelper.CreateSelfSignedSsl(
		pkiConfDir, newHostname, aliasesHostnames,
	)
	if createCertErr != nil {
		return errors.New("CreateSelfSignedSslFailed: " + createCertErr.Error())
	}

	replaceCertPathErr := helpers.replaceVirtualHostSslCertificate(newHostname)
	if replaceCertPathErr != nil {
		return errors.New(
			"ReplaceVirtualHostSslCertificateFailed: " + replaceCertPathErr.Error(),
		)
	}

	return helpers.ReloadWebServer()
}
