package sharedHelper

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
)

var KnownParamConstructors = map[string]func(any) (any, error){
	"username":          tkPresentation.ParamConstructorWrapper(valueObject.NewUsername),
	"password":          tkPresentation.ParamConstructorWrapper(valueObject.NewPassword),
	"operatorIpAddress": tkPresentation.ParamConstructorWrapper(valueObject.NewIpAddress),
	"operatorAccountId": tkPresentation.ParamConstructorWrapper(valueObject.NewAccountId),
}
