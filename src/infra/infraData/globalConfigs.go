package infraData

type globalConfigs struct {
	PkiConfDir                       string
	DomainOwnershipValidationUrlPath string
	PrimaryPublicDir                 string
}

var GlobalConfigs = globalConfigs{
	PkiConfDir:                       "/app/conf/pki",
	DomainOwnershipValidationUrlPath: "/validateOwnership",
	PrimaryPublicDir:                 "/app/html",
}
