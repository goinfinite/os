package infraEnvs

const (
	InfiniteOsVersion            string = "0.1.2"
	InfiniteOsMainDir            string = "/infinite"
	InfiniteOsBinary             string = InfiniteOsMainDir + "/os"
	PersistentDatabaseFilePath   string = InfiniteOsMainDir + "/os.db"
	TrailDatabaseFilePath        string = InfiniteOsMainDir + "/trail.db"
	MarketplaceItemsDir          string = InfiniteOsMainDir + "/marketplace"
	PrimaryPublicDir             string = "/app/html"
	VirtualHostsConfDir          string = "/app/conf/nginx"
	PrimaryVhostFileName         string = "primary.conf"
	MappingsConfDir              string = "/app/conf/nginx/mapping"
	PkiConfDir                   string = "/app/conf/pki"
	PhpWebserverMainConfFilePath string = "/usr/local/lsws/conf/httpd_config.conf"
	AccessTokenCookieKey         string = "os-access-token"
)
