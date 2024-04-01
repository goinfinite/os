package infraShared

type envDataInfra struct {
	PkiConfDir                       string
	DomainOwnershipValidationUrlPath string
}

var EnvDataInfra = envDataInfra{
	PkiConfDir:                       "/app/conf/pki",
	DomainOwnershipValidationUrlPath: "/validateOwnership",
}
