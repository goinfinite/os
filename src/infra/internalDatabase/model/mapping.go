package dbModel

import (
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type Mapping struct {
	ID           uint `gorm:"primarykey"`
	Hostname     string
	Path         string
	MatchPattern string
	TargetType   string
	TargetValue  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
) Mapping {
	mappingModel := Mapping{
		Hostname:     hostname,
		Path:         path,
		MatchPattern: matchPattern,
		TargetType:   targetType,
		TargetValue:  targetValue,
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

	var targetServiceNamePtr *valueObject.ServiceName
	if targetType.String() == "service" {
		targetServiceName, err := valueObject.NewServiceName(*model.TargetValue)
		if err != nil {
			return mapping, err
		}
		targetServiceNamePtr = &targetServiceName
	}

	var targetUrlPtr *valueObject.Url
	if targetType.String() == "url" {
		targetUrl, err := valueObject.NewUrl(*model.TargetValue)
		if err != nil {
			return mapping, err
		}
		targetUrlPtr = &targetUrl
	}

	var targetHttpResponseCodePtr *valueObject.HttpResponseCode
	if targetType.String() == "response-code" {
		targetHttpResponseCode, err := valueObject.NewHttpResponseCode(
			*model.TargetValue,
		)
		if err != nil {
			return mapping, err
		}
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	var targetInlineHtmlContentPtr *valueObject.InlineHtmlContent
	if targetType.String() == "inline-html" {
		targetInlineHtmlContent, err := valueObject.NewInlineHtmlContent(
			*model.TargetValue,
		)
		if err != nil {
			return mapping, err
		}
		targetInlineHtmlContentPtr = &targetInlineHtmlContent
	}

	return entity.NewMapping(
		mappingId,
		hostname,
		path,
		matchPattern,
		targetType,
		targetServiceNamePtr,
		targetUrlPtr,
		targetHttpResponseCodePtr,
		targetInlineHtmlContentPtr,
	), nil
}

func (Mapping) AddDtoToModel(createDto dto.CreateMapping) Mapping {
	var targetValuePtr *string
	switch createDto.TargetType.String() {
	case "url":
		targetUrlStr := createDto.TargetUrl.String()
		targetValuePtr = &targetUrlStr
	case "service":
		targetServiceNameStr := createDto.TargetServiceName.String()
		targetValuePtr = &targetServiceNameStr
	case "response-code":
		targetHttpResponseCodeStr := createDto.TargetHttpResponseCode.String()
		targetValuePtr = &targetHttpResponseCodeStr
	case "inline-html":
		targetInlineHtmlContentStr := createDto.TargetInlineHtmlContent.String()
		targetValuePtr = &targetInlineHtmlContentStr
	}

	return NewMapping(
		0,
		createDto.Hostname.String(),
		createDto.Path.String(),
		createDto.MatchPattern.String(),
		createDto.TargetType.String(),
		targetValuePtr,
	)
}
