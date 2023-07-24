package infra

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type WsQueryRepo struct {
}

func (ws WsQueryRepo) GetVirtualHosts() ([]valueObject.Fqdn, error) {
	olsConfigFile := "/usr/local/lsws/conf/httpd_config.conf"
	output, err := infraHelper.RunCmd(
		"awk",
		"/virtualhost /{print $2}",
		olsConfigFile,
	)
	if err != nil {
		log.Printf("FailedToGetVirtualHosts: %v", err)
		return nil, errors.New("FailedToGetVirtualHosts")
	}

	virtualHosts := []valueObject.Fqdn{
		valueObject.NewFqdnPanic(os.Getenv("VIRTUAL_HOST")),
	}
	for _, vhost := range strings.Split(output, "\n") {
		if vhost == "" || vhost == "app" {
			continue
		}

		virtualHost, err := valueObject.NewFqdn(vhost)
		if err != nil {
			continue
		}

		virtualHosts = append(virtualHosts, virtualHost)
	}

	return virtualHosts, nil
}
