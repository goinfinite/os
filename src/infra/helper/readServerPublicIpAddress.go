package infraHelper

import (
	"errors"
	"strings"

	"github.com/goinfinite/os/src/domain/valueObject"
)

func ReadServerPublicIpAddress() (ipAddress valueObject.IpAddress, err error) {
	digCmd := "dig +short TXT"
	rawRecord, err := RunCmd(RunCmdSettings{
		Command:               digCmd + " o-o.myaddr.l.google.com @ns1.google.com",
		ShouldRunWithSubShell: true,
	})
	if err != nil || rawRecord == "" {
		rawRecord, err = RunCmd(RunCmdSettings{
			Command:               digCmd + " CH whoami.cloudflare @1.1.1.1",
			ShouldRunWithSubShell: true,
		})
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
