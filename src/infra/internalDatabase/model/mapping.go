package dbModel

import (
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type Mapping struct {
	ID                         uint   `gorm:"primarykey"`
	Hostname                   string `gorm:"not null"`
	Path                       string `gorm:"not null"`
	MatchPattern               string `gorm:"not null"`
	TargetType                 string `gorm:"not null"`
	TargetValue                *string
	TargetHttpResponseCode     *string
	VirtualHostHostname        string `gorm:"not null"`
	MarketplaceInstalledItemId *uint
	CreatedAt                  time.Time `gorm:"not null"`
	UpdatedAt                  time.Time `gorm:"not null"`
}

func (Mapping) TableName() string {
	return "mappings"
}

func NewMapping(
	id uint,
	hostname string,
	path string,
	matchPattern string,
	targetType string,
	targetValue *string,
	targetHttpResponseCode *string,
	vhostName string,
) Mapping {
	mappingModel := Mapping{
		Hostname:               hostname,
		Path:                   path,
		MatchPattern:           matchPattern,
		TargetType:             targetType,
		TargetValue:            targetValue,
		TargetHttpResponseCode: targetHttpResponseCode,
		VirtualHostHostname:    vhostName,
	}

	if id != 0 {
		mappingModel.ID = id
	}

	return mappingModel
}

func (model Mapping) ToEntity() (entity.Mapping, error) {
	var mapping entity.Mapping

	mappingId, err := valueObject.NewMappingId(model.ID)
	if err != nil {
		return mapping, err
	}

	hostname, err := valueObject.NewFqdn(model.Hostname)
	if err != nil {
		return mapping, err
	}

	path, err := valueObject.NewMappingPath(model.Path)
	if err != nil {
		return mapping, err
	}

	matchPattern, err := valueObject.NewMappingMatchPattern(model.MatchPattern)
	if err != nil {
		return mapping, err
	}

	targetType, err := valueObject.NewMappingTargetType(model.TargetType)
	if err != nil {
		return mapping, err
	}

	var targetValuePtr *valueObject.MappingTargetValue
	if model.TargetValue != nil {
		targetValue, err := valueObject.NewMappingTargetValue(
			*model.TargetValue, targetType,
		)
		if err != nil {
			return mapping, err
		}
		targetValuePtr = &targetValue
	}

	var targetHttpResponseCodePtr *valueObject.HttpResponseCode
	if model.TargetHttpResponseCode != nil {
		targetHttpResponseCode, err := valueObject.NewHttpResponseCode(
			*model.TargetHttpResponseCode,
		)
		if err != nil {
			return mapping, err
		}
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	return entity.NewMapping(
		mappingId,
		hostname,
		path,
		matchPattern,
		targetType,
		targetValuePtr,
		targetHttpResponseCodePtr,
	), nil
}

func (Mapping) ToModel(mappingEntity entity.Mapping) Mapping {
	var targetValuePtr *string
	if mappingEntity.TargetValue != nil {
		targetValueStr := mappingEntity.TargetValue.String()
		targetValuePtr = &targetValueStr
	}

	var targetHttpResponseCodePtr *string
	if mappingEntity.TargetHttpResponseCode != nil {
		targetHttpResponseCodeStr := mappingEntity.TargetHttpResponseCode.String()
		targetHttpResponseCodePtr = &targetHttpResponseCodeStr
	}

	return NewMapping(
		uint(mappingEntity.Id),
		mappingEntity.Hostname.String(),
		mappingEntity.Path.String(),
		mappingEntity.MatchPattern.String(),
		mappingEntity.TargetType.String(),
		targetValuePtr,
		targetHttpResponseCodePtr,
		"",
	)
}
