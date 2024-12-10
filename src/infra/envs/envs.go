package infraEnvs

const (
	InfiniteOsVersion            string = "0.1.5"
	InfiniteOsMainDir            string = "/infinite"
	InfiniteOsBinary             string = InfiniteOsMainDir + "/os"
	InfiniteOsEnvFilePath        string = InfiniteOsMainDir + "/.env"
	PersistentDatabaseFilePath   string = InfiniteOsMainDir + "/os.db"
	TrailDatabaseFilePath        string = InfiniteOsMainDir + "/trail.db"
	MarketplaceCatalogItemsDir   string = InfiniteOsMainDir + "/marketplace"
	InstallableServicesItemsDir  string = InfiniteOsMainDir + "/services"
	PrimaryPublicDir             string = "/app/html"
	VirtualHostsConfDir          string = "/app/conf/nginx"
	PrimaryVhostFileName         string = "primary.conf"
	MappingsConfDir              string = "/app/conf/nginx/mapping"
	PkiConfDir                   string = "/app/conf/pki"
	PhpWebserverMainConfFilePath string = "/usr/local/lsws/conf/httpd_config.conf"
	AccessTokenCookieKey         string = "os-access-token"
	UserDataBaseDirectory        string = "/home/"
)
