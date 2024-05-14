package dbModel

import (
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"github.com/speedianet/os/src/infra/infraData"
)

type VirtualHost struct {
	ID             uint   `gorm:"primarykey"`
	Hostname       string `gorm:"not null"`
	Type           string `gorm:"not null"`
	RootDirectory  string `gorm:"not null"`
	ParentHostname *string
	Mappings       []Mapping
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
}

func (model VirtualHost) InitialEntries() []interface{} {
	primaryVhost, _ := infraHelper.GetPrimaryVirtualHost()
	primaryEntry := VirtualHost{
		ID:            1,
		Hostname:      primaryVhost.String(),
		Type:          "primary",
		RootDirectory: infraData.GlobalConfigs.PrimaryPublicDir,
	}

	return []interface{}{primaryEntry}
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
	if model.ParentHostname != nil {
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

func (VirtualHost) ToModel(entity entity.VirtualHost) (VirtualHost, error) {
	var parentHostnamePtr *string
	if entity.ParentHostname != nil {
		parentHostnameStr := entity.ParentHostname.String()
		parentHostnamePtr = &parentHostnameStr
	}

	mappings := []Mapping{}
	for _, mapping := range entity.Mappings {
		mappingEntity := Mapping{}.ToModel(mapping)
		mappingEntity.ID = entity.Id.Get()
		mappings = append(mappings, mappingEntity)
	}

	return VirtualHost{
		ID:             entity.Id.Get(),
		Hostname:       entity.Hostname.String(),
		Type:           entity.Type.String(),
		RootDirectory:  entity.RootDirectory.String(),
		ParentHostname: parentHostnamePtr,
		Mappings:       mappings,
	}, nil
}
