package vhostInfra

import (
	"errors"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"

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
	rawPrimaryConfContentStr, err := tkInfra.FileClerk{}.ReadFileContent(
		infraEnvs.PrimaryVirtualHostConfPath, nil,
	)
	if err != nil {
		return primaryHostname, err
	}

	serverNameRegex := regexp.MustCompile(`^\s*server_name\s+([^;]+);`)
	for rawPrimaryConfLineStr := range strings.SplitSeq(
		rawPrimaryConfContentStr, "\n",
	) {
		rawServerNameMatchStrSlice := serverNameRegex.FindStringSubmatch(
			rawPrimaryConfLineStr,
		)
		if len(rawServerNameMatchStrSlice) < 2 {
			continue
		}

		rawServerNameValueStr := rawServerNameMatchStrSlice[1]
		rawPrimaryConfServerNamesStrSlice := strings.Fields(
			rawServerNameValueStr,
		)
		if len(rawPrimaryConfServerNamesStrSlice) == 0 {
			continue
		}

		return tkValueObject.NewFqdn(rawPrimaryConfServerNamesStrSlice[0])
	}

	return primaryHostname, errors.New("PrimaryServerNameNotFound")
}

// ReadPrimaryVirtualHostHostname returns the primary virtual host hostname from the
// environment variables, followed by the web server configuration and fallback
// to the shell command.
// If you need the hostname straight from the web server configuration, use
// ReadPrimaryVirtualHostHostnameFromWebServerConf instead.
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

	reloadOutput, reloadErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: infraEnvs.WebServerBinaryPath,
		Args:    []string{"-s", "reload", "-c", infraEnvs.WebServerMainConfPath},
	}).Run()
	if reloadErr != nil {
		combinedOutput := reloadOutput + " " + reloadErr.Error()
		failBecauseWebServerNotRunning := strings.Contains(combinedOutput, "nginx.pid")
		if !failBecauseWebServerNotRunning {
			return errors.New("WebServerReloadFail: " + combinedOutput)
		}
	}

	time.Sleep(1 * time.Second)

	return nil
}
