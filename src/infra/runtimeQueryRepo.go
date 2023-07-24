package infra

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type RuntimeQueryRepo struct {
}

func (r RuntimeQueryRepo) GetPhpVersions() ([]valueObject.PhpVersion, error) {
	olsConfigFile := "/usr/local/lsws/conf/httpd_config.conf"
	output, err := infraHelper.RunCmd(
		"awk",
		"/extprocessor lsphp/{print $2}",
		olsConfigFile,
	)
	if err != nil {
		log.Printf("FailedToGetPhpVersions: %v", err)
		return nil, errors.New("FailedToGetPhpVersions")
	}

	phpVersions := []valueObject.PhpVersion{}
	for _, version := range strings.Split(output, "\n") {
		if version == "" {
			continue
		}

		version = strings.Replace(version, "lsphp", "", 1)
		phpVersion, err := valueObject.NewPhpVersion(version)
		if err != nil {
			continue
		}

		phpVersions = append(phpVersions, phpVersion)
	}

	return phpVersions, nil
}
