package infraData

type globalConfigs struct {
	PkiConfDir                       string
	DomainOwnershipValidationUrlPath string
	PrimaryPublicDir                 string
	VirtualHostsConfDir              string
	MappingsConfDir                  string
	PublicIpTransientKey             string
	DatabaseFilePath                 string
	PrimaryVhostFileName             string
	OlsHttpdConfFilePath             string
}

var GlobalConfigs = globalConfigs{
	PkiConfDir:                       "/app/conf/pki",
	DomainOwnershipValidationUrlPath: "/validateOwnership",
	PrimaryPublicDir:                 "/app/html",
	VirtualHostsConfDir:              "/app/conf/nginx",
	MappingsConfDir:                  "/app/conf/nginx/mapping",
	PublicIpTransientKey:             "PublicIp",
	DatabaseFilePath:                 "/speedia/os.db",
	PrimaryVhostFileName:             "primary",
	OlsHttpdConfFilePath:             "/usr/local/lsws/conf/httpd_config.conf",
}
