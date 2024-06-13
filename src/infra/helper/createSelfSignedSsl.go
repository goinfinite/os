package infraHelper

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

func altNamesFactory(
	virtualHostHostname string,
	aliasesHostname []string,
) []string {
	virtualHostHostnameWithWww := "www." + virtualHostHostname
	altNamesValues := []string{virtualHostHostname, virtualHostHostnameWithWww}
	for _, aliasHostname := range aliasesHostname {
		aliasHostnameWithWww := "www." + aliasHostname
		altNamesValues = append(altNamesValues, aliasHostname, aliasHostnameWithWww)
	}

	altNamesList := []string{}
	for altNameIndex, altName := range altNamesValues {
		dnsIndex := strconv.Itoa(altNameIndex)
		dnsAltNameConf := "DNS." + dnsIndex + " = " + altName

		altNamesList = append(altNamesList, dnsAltNameConf)
	}

	return altNamesList
}

func selfSignedConfFileFactory(
	virtualHostHostname string,
	aliasesHostname []string,
) (string, error) {
	altNames := altNamesFactory(virtualHostHostname, aliasesHostname)
	valuesToInterpolate := map[string]interface{}{
		"VirtualHostHostname": virtualHostHostname,
		"AltNames":            altNames,
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
CN = {{ .VirtualHostHostname }}

[ v3_req ]
subjectAltName = @alt_names

[ alt_names ]
{{- range $altName := .AltNames }}
{{ $altName }}
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
	virtualHostHostname string,
	aliasesHostname []string,
) error {
	selfSignedConfContent, err := selfSignedConfFileFactory(
		virtualHostHostname, aliasesHostname,
	)
	if err != nil {
		return errors.New("GenerateSelfSignedConfFileError: " + err.Error())
	}
	log.Print(selfSignedConfContent)

	selfSignedConfTempFilePath := "/tmp/" + virtualHostHostname + "_selfSignedSsl.conf"
	shouldOverwrite := true
	err = UpdateFile(selfSignedConfTempFilePath, selfSignedConfContent, shouldOverwrite)
	if err != nil {
		return errors.New("GenerateSelfSignedConfFileError: " + err.Error())
	}

	vhostCertKeyFilePath := dirPath + "/" + virtualHostHostname + ".key"
	vhostCertFilePath := dirPath + "/" + virtualHostHostname + ".crt"

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
			"CreateSelfSignedSslFailed (" + virtualHostHostname + "): " + err.Error(),
		)
	}

	err = os.Remove(selfSignedConfTempFilePath)
	if err != nil {
		return errors.New("DeleteSelfSignedConfFileError: " + err.Error())
	}

	return nil
}
