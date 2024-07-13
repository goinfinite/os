package service

type StatusEnum string

const (
	Success      StatusEnum = "success"
	Created      StatusEnum = "created"
	MultiStatus  StatusEnum = "multiStatus"
	UserError    StatusEnum = "userError"
	Unauthorized StatusEnum = "unauthorized"
	InfraError   StatusEnum = "infraError"
)

type ServiceOutput struct {
	Status StatusEnum  `json:"status"`
	Body   interface{} `json:"body"`
}

func NewServiceOutput(status StatusEnum, body interface{}) ServiceOutput {
	return ServiceOutput{
		Status: status,
		Body:   body,
	}
}
