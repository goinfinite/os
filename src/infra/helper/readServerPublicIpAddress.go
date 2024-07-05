package infraHelper

import (
	"errors"
	"strings"

	"github.com/speedianet/os/src/domain/valueObject"
)

func ReadServerPublicIpAddress() (ipAddress valueObject.IpAddress, err error) {
	digCmd := "dig +short TXT"
	rawRecord, err := RunCmdWithSubShell(
		digCmd + " o-o.myaddr.l.google.com @ns1.google.com",
	)
	if err != nil || rawRecord == "" {
		rawRecord, err = RunCmdWithSubShell(
			digCmd + "CH whoami.cloudflare @1.1.1.1",
		)
		if err != nil {
			return ipAddress, err
		}
	}

	rawRecord = strings.Trim(rawRecord, `"`)
	rawRecord = strings.TrimSpace(rawRecord)
	if rawRecord == "" {
		return ipAddress, errors.New("PublicIpAddressNotFound")
	}

	ipAddress, err = valueObject.NewIpAddress(rawRecord)
	if err != nil {
		return ipAddress, err
	}

	return ipAddress, nil
}
