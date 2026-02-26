package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type RunPhpCommandRequest struct {
	Hostname          tkValueObject.Fqdn        `json:"hostname"`
	Command           tkValueObject.UnixCommand `json:"command"`
	TimeoutSecs       *uint64                   `json:"timeoutSecs"`
	OperatorAccountId tkValueObject.AccountId   `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress   `json:"-"`
}

func NewRunPhpCommandRequest(
	hostname tkValueObject.Fqdn,
	command tkValueObject.UnixCommand,
	timeoutSecs *uint64,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) RunPhpCommandRequest {
	return RunPhpCommandRequest{
		Hostname:          hostname,
		Command:           command,
		TimeoutSecs:       timeoutSecs,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}

type RunPhpCommandResponse struct {
	StdOutput *valueObject.UnixCommandOutput `json:"stdOut"`
	StdError  *valueObject.UnixCommandOutput `json:"stdErr"`
	ExitCode  *int                           `json:"exitCode"`
}
