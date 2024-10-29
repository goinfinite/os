package infraEnvs

const (
	InfiniteOsVersion            string = "0.1.2"
	InfiniteOsBinary             string = "/infinite/os"
	PersistentDatabaseFilePath   string = "/infinite/os.db"
	TrailDatabaseFilePath        string = "/infinite/trail.db"
	PrimaryPublicDir             string = "/app/html"
	VirtualHostsConfDir          string = "/app/conf/nginx"
	PrimaryVhostFileName         string = "primary.conf"
	MappingsConfDir              string = "/app/conf/nginx/mapping"
	PkiConfDir                   string = "/app/conf/pki"
	PhpWebserverMainConfFilePath string = "/usr/local/lsws/conf/httpd_config.conf"
	AccessTokenCookieKey         string = "os-access-token"
)
