package infra

import (
	"github.com/speedianet/os/src/domain/dto"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
)

type ActivityRecordCmdRepo struct {
	trailDbSvc *internalDbInfra.TrailDatabaseService
}

func NewActivityRecordCmdRepo(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *ActivityRecordCmdRepo {
	return &ActivityRecordCmdRepo{trailDbSvc: trailDbSvc}
}

func (repo *ActivityRecordCmdRepo) Create(createDto dto.CreateActivityRecord) error {
	var codePtr *string
	if createDto.Code != nil {
		code := createDto.Code.String()
		codePtr = &code
	}

	var messagePtr *string
	if createDto.Message != nil {
		message := createDto.Message.String()
		messagePtr = &message
	}

	var ipAddressPtr *string
	if createDto.IpAddress != nil {
		ipAddress := createDto.IpAddress.String()
		ipAddressPtr = &ipAddress
	}

	var operatorAccountIdPtr *uint64
	if createDto.OperatorAccountId != nil {
		operatorAccountId := createDto.OperatorAccountId.Read()
		operatorAccountIdPtr = &operatorAccountId
	}

	var targetAccountIdPtr *uint64
	if createDto.TargetAccountId != nil {
		targetAccountId := createDto.TargetAccountId.Read()
		targetAccountIdPtr = &targetAccountId
	}

	var usernamePtr *string
	if createDto.Username != nil {
		username := createDto.Username.String()
		usernamePtr = &username
	}

	var mappingIdPtr *uint64
	if createDto.MappingId != nil {
		mappingId := createDto.MappingId.Get()
		mappingIdPtr = &mappingId
	}

	securityEventModel := dbModel.NewActivityRecord(
		0, createDto.Level.String(), codePtr, messagePtr, ipAddressPtr,
		operatorAccountIdPtr, targetAccountIdPtr, usernamePtr, mappingIdPtr,
	)

	return repo.trailDbSvc.Handler.Create(&securityEventModel).Error
}
