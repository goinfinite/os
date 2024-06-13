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
	vhostNameWithWww := "www." + vhostName
	altNames := []string{vhostName, vhostNameWithWww}
	for _, aliasHostname := range aliasesHostname {
		aliasHostnameWithWww := "www." + aliasHostname
		altNames = append(altNames, aliasHostname, aliasHostnameWithWww)
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

	selfSignedConfFileTemplate := `[ req ]
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

	selfSignedConfFileTemplatePtr, err := template.
		New("selfSignedConfFile").
		Parse(selfSignedConfFileTemplate)
	if err != nil {
		return "", errors.New("TemplateParsingError: " + err.Error())
	}

	var selfSignedConfFileContent strings.Builder
	err = selfSignedConfFileTemplatePtr.Execute(
		&selfSignedConfFileContent,
		valuesToInterpolate,
	)
	if err != nil {
		return "", errors.New("TemplateExecutionError: " + err.Error())
	}

	return selfSignedConfFileContent.String(), nil
}

func CreateSelfSignedSsl(
	dirPath string,
	vhostName string,
	aliasesHostname []string,
) error {
	selfSignedConfContent, err := selfSignedConfFileFactory(
		vhostName, aliasesHostname,
	)
	if err != nil {
		return errors.New("GenerateSelfSignedConfFileError: " + err.Error())
	}

	selfSignedConfTempFilePath := "/tmp/" + vhostName + "_selfSignedSsl.conf"
	shouldOverwrite := true
	err = UpdateFile(selfSignedConfTempFilePath, selfSignedConfContent, shouldOverwrite)
	if err != nil {
		return errors.New("GenerateSelfSignedConfFileError: " + err.Error())
	}

	vhostCertKeyFilePath := dirPath + "/" + vhostName + ".key"
	vhostCertFilePath := dirPath + "/" + vhostName + ".crt"

	_, err = RunCmd(
		"openssl",
		"req",
		"-x509",
		"-nodes",
		"-days",
		"365",
		"-newkey",
		"rsa:2048",
		"-keyout",
		vhostCertKeyFilePath,
		"-out",
		vhostCertFilePath,
		"-config",
		selfSignedConfTempFilePath,
	)
	if err != nil {
		return errors.New(
			"CreateSelfSignedSslFailed (" + vhostName + "): " + err.Error(),
		)
	}

	err = os.Remove(selfSignedConfTempFilePath)
	if err != nil {
		return errors.New("DeleteSelfSignedConfFileError: " + err.Error())
	}

	return nil
}
