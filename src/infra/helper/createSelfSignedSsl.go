package infraHelper

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func altNamesConfFactory(
	vhostName string,
	aliasesHostname []string,
) []string {
	altNames := []string{vhostName, "www." + vhostName}
	for _, aliasHostname := range aliasesHostname {
		altNames = append(altNames, aliasHostname, "www."+aliasHostname)
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
	vhostName string,
	aliasesHostname []string,
) (string, error) {
	altNamesConf := altNamesConfFactory(vhostName, aliasesHostname)
	valuesToInterpolate := map[string]interface{}{
		"VhostName":    vhostName,
		"AltNamesConf": altNamesConf,
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
CN = {{ .VhostName }}

[ v3_req ]
subjectAltName = @alt_names

[ alt_names ]
{{- range $altNameConf := .AltNamesConf }}
{{ $altNameConf }}
{{- end }}
`

	confFileTemplatePtr, err := template.
		New("selfSignedConfFile").
		Parse(confFileTemplate)
	if err != nil {
		return "", errors.New("TemplateParsingError: " + err.Error())
	}

	var confFileContent strings.Builder
	err = confFileTemplatePtr.Execute(
		&confFileContent,
		valuesToInterpolate,
	)
	if err != nil {
		return "", errors.New("TemplateExecutionError: " + err.Error())
	}

	return confFileContent.String(), nil
}

func CreateSelfSignedSsl(
	dirPath string,
	vhostName string,
	aliasesHostname []string,
) error {
	confContent, err := selfSignedConfFileFactory(vhostName, aliasesHostname)
	if err != nil {
		return errors.New("GenerateSelfSignedConfFileError: " + err.Error())
	}

	confTempFilePath := "/tmp/" + vhostName + "_selfSignedSsl.conf"
	shouldOverwrite := true
	err = UpdateFile(confTempFilePath, confContent, shouldOverwrite)
	if err != nil {
		return errors.New("GenerateSelfSignedConfFileError: " + err.Error())
	}

	vhostCertKeyFilePath := dirPath + "/" + vhostName + ".key"
	vhostCertFilePath := dirPath + "/" + vhostName + ".crt"

	_, err = RunCmdWithSubShell(
		"openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout " +
			vhostCertKeyFilePath + " -out " + vhostCertFilePath + " -config " + confTempFilePath,
	)
	if err != nil {
		return errors.New(
			"CreateSelfSignedSslFailed (" + vhostName + "): " + err.Error(),
		)
	}

	err = os.Remove(confTempFilePath)
	if err != nil {
		return errors.New("DeleteSelfSignedConfFileError: " + err.Error())
	}

	return nil
}
