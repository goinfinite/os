package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type RunPhpCommandRequest struct {
	Hostname          valueObject.Fqdn        `json:"hostname"`
	Command           valueObject.UnixCommand `json:"command"`
	TimeoutSecs       *uint64                 `json:"timeoutSecs"`
	OperatorAccountId valueObject.AccountId   `json:"-"`
	OperatorIpAddress valueObject.IpAddress   `json:"-"`
}

func NewRunPhpCommandRequest(
	hostname valueObject.Fqdn,
	command valueObject.UnixCommand,
	timeoutSecs *uint64,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
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
