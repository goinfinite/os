package dbModel

import (
	"errors"
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type VirtualHost struct {
	Hostname       string `gorm:"primarykey;not null"`
	Type           string `gorm:"not null"`
	RootDirectory  string `gorm:"not null"`
	ParentHostname *string
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

func (model VirtualHost) InitialEntries() (entries []interface{}, err error) {
	primaryVhostName, err := infraHelper.ReadPrimaryVirtualHostHostname()
	if err != nil {
		return entries, errors.New("ReadPrimaryVirtualHostHostnameError: " + err.Error())
	}

	primaryEntry := VirtualHost{
		Hostname:      primaryVhostName.String(),
		Type:          "primary",
		RootDirectory: infraEnvs.PrimaryPublicDir,
	}

	return []interface{}{primaryEntry}, nil
}

func (model VirtualHost) ToEntity() (vhost entity.VirtualHost, err error) {
	hostname, err := valueObject.NewFqdn(model.Hostname)
	if err != nil {
		return vhost, err
	}

	vhostType, err := valueObject.NewVirtualHostType(model.Type)
	if err != nil {
		return vhost, err
	}

	rootDir, err := valueObject.NewUnixFilePath(model.RootDirectory)
	if err != nil {
		return vhost, err
	}

	var parentHostnamePtr *valueObject.Fqdn
	if model.ParentHostname != nil {
		parentHostname, err := valueObject.NewFqdn(*model.ParentHostname)
		if err != nil {
			return vhost, err
		}
		parentHostnamePtr = &parentHostname
	}

	return entity.NewVirtualHost(hostname, vhostType, rootDir, parentHostnamePtr), nil
}
