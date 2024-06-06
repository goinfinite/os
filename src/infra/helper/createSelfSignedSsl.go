package infraHelper

import (
	"errors"
	"os"
	"strings"
	"text/template"
)

func selfSignedConfFileFactory(
	virtualHost string,
	aliases []string,
) (string, error) {
	altNames := []string{virtualHost}
	altNames = append(altNames, aliases...)

	valuesToInterpolate := map[string]interface{}{
		"VirtualHost": virtualHost,
		"AltNames":    altNames,
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
CN = {{ .VirtualHost }}

[ v3_req ]
subjectAltName = @alt_names

[ alt_names ]
{{- $lastDnsIndex := 0 }}
{{- range $altName := .AltNames }}
{{- $dnsIndex := increaseIndex $lastDnsIndex }}
DNS.{{ $dnsIndex }} = {{ $altName }}
{{- $wwwDnsIndex := increaseIndex $dnsIndex }}
DNS.{{ $wwwDnsIndex }} = www.{{ $altName }}
{{- $lastDnsIndex = $wwwDnsIndex }}
{{- end }}
`

	selfSignedConfFileTemplatePtr := template.New("selfSignedConfFile")
	selfSignedConfFileTemplatePtr = selfSignedConfFileTemplatePtr.Funcs(
		template.FuncMap{
			"increaseIndex": func(currentIndex int) int {
				return currentIndex + 1
			},
		},
	)

	selfSignedConfFileTemplatePtr, err := selfSignedConfFileTemplatePtr.Parse(selfSignedConfFileTemplate)
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
	virtualHost string,
	aliases []string,
) error {
	selfSignedConfContent, err := selfSignedConfFileFactory(virtualHost, aliases)
	if err != nil {
		return errors.New("GenerateSelfSignedConfFileError: " + err.Error())
	}

	selfSignedConfFilePath := "/tmp/" + virtualHost + "_selfSignedSsl.conf"
	shouldOverwrite := true
	err = UpdateFile(selfSignedConfFilePath, selfSignedConfContent, shouldOverwrite)
	if err != nil {
		return errors.New("GenerateSelfSignedConfFileError: " + err.Error())
	}

	vhostCertKeyFilePath := dirPath + "/" + virtualHost + ".key"
	vhostCertFilePath := dirPath + "/" + virtualHost + ".crt"

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
		selfSignedConfFilePath,
	)
	if err != nil {
		return errors.New("CreateSelfSignedSslFailed (" + virtualHost + "): " + err.Error())
	}

	err = os.Remove(selfSignedConfFilePath)
	if err != nil {
		return errors.New("DeleteSelfSignedConfFileError: " + err.Error())
	}

	return nil
}
