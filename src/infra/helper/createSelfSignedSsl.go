package infraHelper

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/goinfinite/os/src/domain/valueObject"
)

func altNamesConfFactory(
	vhostHostname valueObject.Fqdn,
	aliasesHostname []valueObject.Fqdn,
) []string {
	altNames := []string{vhostHostname.String(), "www." + vhostHostname.String()}
	for _, aliasHostname := range aliasesHostname {
		altNames = append(altNames, aliasHostname.String(), "www."+aliasHostname.String())
	}

	altNamesConfList := []string{}
	for altNameIndex, altName := range altNames {
		dnsIndex := strconv.Itoa(altNameIndex)
		altNameConf := "DNS." + dnsIndex + " = " + altName

		altNamesConfList = append(altNamesConfList, altNameConf)
	}

	return altNamesConfList
}

func selfSignedConfFileFactory(
	vhostHostname valueObject.Fqdn,
	aliasesHostname []valueObject.Fqdn,
) (string, error) {
	altNamesConf := altNamesConfFactory(vhostHostname, aliasesHostname)
	valuesToInterpolate := map[string]interface{}{
		"VirtualHostHostname": vhostHostname,
		"AltNamesConf":        altNamesConf,
	}

	confFileTemplate := `[ req ]
default_bits = 2048
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[ req_distinguished_name ]
C = US
ST = California
L = Los Angeles
CN = {{ .VirtualHostHostname }}

[ v3_req ]
subjectAltName = @alt_names

[ alt_names ]
{{- range $altNameConf := .AltNamesConf }}
{{ $altNameConf }}
{{- end }}
`

	confFileTemplatePtr, err := template.New("selfSignedConfFile").Parse(confFileTemplate)
	if err != nil {
		return "", errors.New("TemplateParsingError: " + err.Error())
	}

	var confFileContent strings.Builder
	err = confFileTemplatePtr.Execute(&confFileContent, valuesToInterpolate)
	if err != nil {
		return "", errors.New("TemplateExecutionError: " + err.Error())
	}

	return confFileContent.String(), nil
}

func CreateSelfSignedSsl(
	dirPath valueObject.UnixFilePath,
	vhostHostname valueObject.Fqdn,
	aliasesHostname []valueObject.Fqdn,
) error {
	confContent, err := selfSignedConfFileFactory(vhostHostname, aliasesHostname)
	if err != nil {
		return errors.New("SelfSignedConfFactoryError: " + err.Error())
	}

	vhostHostnameStr := vhostHostname.String()
	confTempFilePath := "/tmp/" + vhostHostnameStr + "_selfSignedSsl.conf"
	shouldOverwrite := true
	err = UpdateFile(confTempFilePath, confContent, shouldOverwrite)
	if err != nil {
		return errors.New("UpdateSelfSignedConfFileError: " + err.Error())
	}

	dirPathStr := dirPath.String()
	vhostCertKeyFilePath := dirPathStr + "/" + vhostHostnameStr + ".key"
	vhostCertFilePath := dirPathStr + "/" + vhostHostnameStr + ".crt"

	createSslCmd := "openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout " +
		vhostCertKeyFilePath + " -out " + vhostCertFilePath + " -config " + confTempFilePath
	_, err = RunCmd(RunCmdSettings{
		Command:               createSslCmd,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return errors.New("CreateSelfSignedSslCmdFailed: " + err.Error())
	}

	return os.Remove(confTempFilePath)
}
