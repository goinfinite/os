package dbModel

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type VirtualHost struct {
	ID             uint   `gorm:"primarykey"`
	Hostname       string `gorm:"not null"`
	Type           string `gorm:"not null"`
	RootDirectory  string `gorm:"not null"`
	ParentHostname *string
	Mappings       []Mapping
}

func (model VirtualHost) ToEntity() (vhost entity.VirtualHost, err error) {
	id, err := valueObject.NewVirtualHostId(model.ID)
	if err != nil {
		return vhost, err
	}

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
	if model.ParentHostname == nil {
		parentHostname, err := valueObject.NewFqdn(*model.ParentHostname)
		if err != nil {
			return vhost, err
		}
		parentHostnamePtr = &parentHostname
	}

	mappings := []entity.Mapping{}
	if len(model.Mappings) > 0 {
		for _, mappingModel := range model.Mappings {
			mapping, err := mappingModel.ToEntity()
			if err != nil {
				return vhost, err
			}
			mappings = append(mappings, mapping)
		}
	}

	return entity.NewVirtualHost(
		id,
		hostname,
		vhostType,
		rootDir,
		parentHostnamePtr,
		mappings,
	), nil
}
