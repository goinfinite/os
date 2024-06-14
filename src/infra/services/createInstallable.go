package servicesInfra

import (
	"embed"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"github.com/speedianet/os/src/infra/infraData"
)

var supportedServicesVersion = map[string]string{
	"mariadb":    `^(10\.([6-9]|10|11)|11\.[0-9]{1,2})$`,
	"postgresql": `^(1[2-6])$`,
	"node":       `^(1[2-9]|2[0-2]])$`,
	"redis":      `^6\.(0|2)|^7\.0$`,
}

var OlsPackages = []string{
	"openlitespeed",
}

var PhpPackages = []string{
	"lsphp74",
	"lsphp74-common",
	"lsphp74-curl",
	"lsphp74-intl",
	"lsphp74-mysql",
	"lsphp74-opcache",
	"lsphp80",
	"lsphp80-common",
	"lsphp80-curl",
	"lsphp80-intl",
	"lsphp80-mysql",
	"lsphp80-opcache",
	"lsphp81",
	"lsphp81-common",
	"lsphp81-curl",
	"lsphp81-intl",
	"lsphp81-mysql",
	"lsphp81-opcache",
	"lsphp82",
	"lsphp82-common",
	"lsphp82-curl",
	"lsphp82-intl",
	"lsphp82-mysql",
	"lsphp82-opcache",
}

var MariaDbPackages = []string{
	"mariadb-server",
}

var RedisPackages = []string{
	"redis-server",
}

//go:embed assets/*
var assets embed.FS

const configurationDir string = "/app/conf"

func copyAssets(srcPath string, dstPath string) error {
	srcPath = "assets/" + srcPath
	srcFile, err := assets.Open(srcPath)
	if err != nil {
		return errors.New("OpenSourceFileError: " + err.Error())
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("OpenDestinationFileError: " + err.Error())
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.New("CopyFileError: " + err.Error())
	}

	return nil
}

func installGpgKey(serviceName string, url string) error {
	keyTempPath := "/speedia/" + serviceName + ".gpg"

	err := infraHelper.DownloadFile(
		url,
		keyTempPath,
	)
	if err != nil {
		return errors.New("DownloadRepoFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"gpg",
		"--batch",
		"--yes",
		"--dearmor",
		"-o",
		"/usr/share/keyrings/"+serviceName+"-archive-keyring.gpg",
		keyTempPath,
	)
	if err != nil {
		return errors.New("GpgImportError: " + err.Error())
	}

	err = os.Remove(keyTempPath)
	if err != nil {
		return errors.New("RemoveRepoFileError: " + err.Error())
	}

	return nil
}

func addPhp() error {
	repoFilePath := "/speedia/repo.litespeed.sh"

	err := infraHelper.DownloadFile(
		"https://repo.litespeed.sh",
		repoFilePath,
	)
	if err != nil {
		return errors.New("DownloadRepoFileError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"bash",
		repoFilePath,
	)
	if err != nil {
		return errors.New("RepoAddError: " + err.Error())
	}

	err = os.Remove(repoFilePath)
	if err != nil {
		return errors.New("RemoveRepoFileError: " + err.Error())
	}

	err = infraHelper.InstallPkgs(OlsPackages)
	if err != nil {
		return errors.New("InstallPhpWebServerError: " + err.Error())
	}

	err = infraHelper.InstallPkgs(PhpPackages)
	if err != nil {
		return err
	}

	os.Symlink(
		"/usr/local/lsws/lsphp82/bin/php",
		"/usr/bin/php",
	)

	err = copyAssets(
		"php/httpd_config.conf",
		"/usr/local/lsws/conf/httpd_config.conf",
	)
	if err != nil {
		return errors.New("CopyAssetsError: " + err.Error())
	}

	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return errors.New("PrimaryVirtualHostNotFound")
	}

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/speedia.net/"+primaryVhost.String()+"/g",
		"/usr/local/lsws/conf/httpd_config.conf",
	)
	if err != nil {
		return errors.New("RenameHttpdVHostError: " + err.Error())
	}

	err = infraHelper.MakeDir(configurationDir + "//php")
	if err != nil {
		return errors.New("CreateConfDirError: " + err.Error())
	}

	err = copyAssets(
		"php/primary.conf",
		configurationDir+"//php/template",
	)
	if err != nil {
		return errors.New("CopyAssetsError: " + err.Error())
	}

	err = copyAssets(
		"php/primary.conf",
		configurationDir+"//php/primary.conf",
	)
	if err != nil {
		return errors.New("CopyAssetsError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/speedia.net/"+primaryVhost.String()+"/g",
		configurationDir+"//php/primary.conf",
	)
	if err != nil {
		return errors.New("RenameVHostError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown",
		"-R",
		"lsadm:nogroup",
		configurationDir+"//php",
	)
	if err != nil {
		return errors.New("ChownConfDirError: " + err.Error())
	}

	err = infraHelper.MakeDir("/app/logs/php-webserver")
	if err != nil {
		return errors.New("CreateLogDirError: " + err.Error())
	}

	chownRecursively := true
	chownSymlinksToo := false
	err = infraHelper.UpdatePermissionsForWebServerUse(
		"/app/logs/php-webserver", chownRecursively, chownSymlinksToo,
	)
	if err != nil {
		return errors.New("ChownLogDirError: " + err.Error())
	}

	httpPortBinding := valueObject.NewPortBinding(
		valueObject.NewNetworkPortPanic(8080),
		valueObject.NewNetworkProtocolPanic("http"),
	)
	httpsPortBinding := valueObject.NewPortBinding(
		valueObject.NewNetworkPortPanic(8443),
		valueObject.NewNetworkProtocolPanic("https"),
	)
	portBindings := []valueObject.PortBinding{
		httpPortBinding,
		httpsPortBinding,
	}

	err = SupervisordFacade{}.CreateConf(
		valueObject.NewServiceNamePanic("php-webserver"),
		valueObject.NewServiceNaturePanic("solo"),
		valueObject.NewServiceTypePanic("runtime"),
		valueObject.NewServiceVersionPanic("latest"),
		valueObject.NewUnixCommandPanic("/usr/local/lsws/bin/litespeed -d"),
		nil,
		portBindings,
		nil,
	)
	if err != nil {
		return errors.New("CreateSupervisorConfError: " + err.Error())
	}

	return nil
}

func addNode(createDto dto.CreateInstallableService) error {
	versionStr := "lts"
	if createDto.Version != nil {
		versionStr = createDto.Version.String()
		re := regexp.MustCompile(supportedServicesVersion["node"])
		isVersionAllowed := re.MatchString(versionStr)

		if !isVersionAllowed {
			return errors.New("InvalidNodeVersion: " + versionStr)
		}
	}

	_, err := infraHelper.RunCmdWithSubShell(
		"mise install node@" + versionStr,
	)
	if err != nil {
		return errors.New("InstallNodeError: " + err.Error())
	}

	err = infraHelper.MakeDir(infraData.GlobalConfigs.PrimaryPublicDir)
	if err != nil {
		return errors.New("CreateBaseDirError: " + err.Error())
	}

	startupFile := valueObject.NewUnixFilePathPanic(
		infraData.GlobalConfigs.PrimaryPublicDir + "/index.js",
	)
	if createDto.StartupFile != nil {
		startupFile = *createDto.StartupFile
	}

	if !infraHelper.FileExists(startupFile.String()) {
		err = copyAssets(
			"nodejs/base-index.js",
			startupFile.String(),
		)
		if err != nil {
			return errors.New("CopyAssetsError: " + err.Error())
		}

		chownRecursively := false
		chownSymlinksToo := false
		err = infraHelper.UpdatePermissionsForWebServerUse(
			startupFile.String(),
			chownRecursively,
			chownSymlinksToo,
		)
		if err != nil {
			return errors.New("ChownDummyIndexError: " + err.Error())
		}
	}

	portBindings := []valueObject.PortBinding{
		valueObject.NewPortBinding(
			valueObject.NewNetworkPortPanic(3000),
			valueObject.NewNetworkProtocolPanic("http"),
		),
	}
	if len(createDto.PortBindings) > 0 {
		portBindings = createDto.PortBindings
	}

	err = SupervisordFacade{}.CreateConf(
		createDto.Name,
		valueObject.NewServiceNaturePanic("multi"),
		valueObject.NewServiceTypePanic("runtime"),
		valueObject.NewServiceVersionPanic(versionStr),
		valueObject.NewUnixCommandPanic(
			"mise x node@"+versionStr+" -- node "+startupFile.String()+" &",
		),
		&startupFile,
		portBindings,
		nil,
	)
	if err != nil {
		return errors.New("CreateSupervisorConfError")
	}

	return nil
}

func addMariaDb(createDto dto.CreateInstallableService) error {
	repoFilePath := "/speedia/repo.mariadb.sh"

	err := infraHelper.DownloadFile(
		"https://r.mariadb.com/downloads/mariadb_repo_setup",
		repoFilePath,
	)
	if err != nil {
		return errors.New("DownloadRepoFileError: " + err.Error())
	}

	versionFlag := ""
	versionStr := "latest"
	if createDto.Version != nil {
		versionStr = createDto.Version.String()
		re := regexp.MustCompile(supportedServicesVersion["mariadb"])
		isVersionAllowed := re.MatchString(versionStr)

		if !isVersionAllowed {
			return errors.New("InvalidMysqlVersion: " + versionStr)
		}

		versionFlag = "--mariadb-server-version=" + versionStr
	}

	_, err = infraHelper.RunCmd(
		"bash",
		repoFilePath,
		versionFlag,
	)
	if err != nil {
		return errors.New("RepoAddError: " + err.Error())
	}

	err = os.Remove(repoFilePath)
	if err != nil {
		return errors.New("RemoveRepoFileError: " + err.Error())
	}

	err = infraHelper.InstallPkgs(MariaDbPackages)
	if err != nil {
		return errors.New("InstallServiceError: " + err.Error())
	}

	os.Symlink("/usr/bin/mariadb", "/usr/bin/mysql")
	os.Symlink("/usr/bin/mariadb-admin", "/usr/bin/mysqladmin")
	os.Symlink("/usr/bin/mariadbd-safe", "/usr/bin/mysqld_safe")

	_, err = infraHelper.RunCmd(
		"mariadbd-safe",
		"--no-watch",
	)
	if err != nil {
		return errors.New("StartMysqldSafeError: " + err.Error())
	}

	time.Sleep(5 * time.Second)

	rootPass := infraHelper.GenPass(16)
	postInstallQueries := []string{
		"ALTER USER 'root'@'localhost' IDENTIFIED BY '" + rootPass + "';",
		"DELETE FROM mysql.user WHERE User='';",
		"DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');",
		"DROP DATABASE IF EXISTS test;",
		"DELETE FROM mysql.db WHERE Db='test' OR Db='test\\_%';",
		"FLUSH PRIVILEGES;",
	}
	postInstallQueriesJoined := strings.Join(postInstallQueries, "; ")
	_, err = infraHelper.RunCmd(
		"mariadb",
		"-e",
		postInstallQueriesJoined,
	)
	if err != nil {
		return errors.New("PostInstallQueryError: " + err.Error())
	}

	err = infraHelper.UpdateFile(
		"/root/.my.cnf",
		"[client]\nuser=root\npassword="+rootPass+"\n",
		true,
	)
	if err != nil {
		return errors.New("CreateMyCnfError: " + err.Error())
	}

	err = os.Chmod("/root/.my.cnf", 0400)
	if err != nil {
		return errors.New("ChmodMyCnfError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"mariadb-admin",
		"--user=root",
		"--password="+rootPass,
		"shutdown",
	)
	if err != nil {
		return errors.New("StopMysqldSafeError: " + err.Error())
	}

	portBindings := []valueObject.PortBinding{
		valueObject.NewPortBinding(
			valueObject.NewNetworkPortPanic(3306),
			valueObject.NewNetworkProtocolPanic("tcp"),
		),
	}

	err = SupervisordFacade{}.CreateConf(
		createDto.Name,
		valueObject.NewServiceNaturePanic("solo"),
		valueObject.NewServiceTypePanic("database"),
		valueObject.NewServiceVersionPanic(versionStr),
		valueObject.NewUnixCommandPanic("/usr/bin/mariadbd-safe"),
		nil,
		portBindings,
		nil,
	)
	if err != nil {
		return errors.New("CreateSupervisorConfError")
	}

	return nil
}

func addPostgresqlDb(createDto dto.CreateInstallableService) error {
	versionStr := "16"
	if createDto.Version != nil {
		versionStr = createDto.Version.String()
		re := regexp.MustCompile(supportedServicesVersion["postgresql"])
		isVersionAllowed := re.MatchString(versionStr)

		if !isVersionAllowed {
			return errors.New("InvalidPostgresqlVersion: " + versionStr)
		}
	}

	err := installGpgKey("postgresql", "https://www.postgresql.org/media/keys/ACCC4CF8.asc")
	if err != nil {
		return errors.New("InstallGpgKeyError: " + err.Error())
	}

	osRelease, err := infraHelper.GetOsRelease()
	if err != nil {
		return errors.New("GetOsReleaseError: " + err.Error())
	}

	repoLine := "deb [signed-by=/usr/share/keyrings/postgresql-archive-keyring.gpg] http://apt.postgresql.org/pub/repos/apt " + osRelease + "-pgdg main"
	err = infraHelper.UpdateFile(
		"/etc/apt/sources.list.d/pgdg.list",
		repoLine,
		true,
	)
	if err != nil {
		return errors.New("CreateRepoFileError: " + err.Error())
	}

	err = infraHelper.InstallPkgs(
		[]string{"postgresql-" + versionStr},
	)
	if err != nil {
		return errors.New("InstallServiceError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"gpasswd",
		"-a",
		"postgres",
		"ssl-cert",
	)
	if err != nil {
		return errors.New("AddPostgresToSslCertError: " + err.Error())
	}

	err = os.Chmod("/etc/ssl/private", 0755)
	if err != nil {
		return errors.New("ChmodSslPrivateError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown",
		"postgres:ssl-cert",
		"/etc/ssl/private/ssl-cert-snakeoil.key",
	)
	if err != nil {
		return errors.New("ChownSslPrivateCertError: " + err.Error())
	}

	err = os.Chmod("/etc/ssl/private/ssl-cert-snakeoil.key", 0600)
	if err != nil {
		return errors.New("ChmodSslPrivateCertError: " + err.Error())
	}

	portBindings := []valueObject.PortBinding{
		valueObject.NewPortBinding(
			valueObject.NewNetworkPortPanic(5432),
			valueObject.NewNetworkProtocolPanic("tcp"),
		),
	}

	postgresUser := valueObject.NewUsernamePanic("postgres")

	err = SupervisordFacade{}.CreateConf(
		createDto.Name,
		valueObject.NewServiceNaturePanic("solo"),
		valueObject.NewServiceTypePanic("database"),
		valueObject.NewServiceVersionPanic(versionStr),
		valueObject.NewUnixCommandPanic(
			"/usr/lib/postgresql/"+versionStr+"/bin/postgres "+
				"-D /var/lib/postgresql/"+versionStr+"/main "+
				"-c config_file=/etc/postgresql/"+versionStr+"/main/postgresql.conf",
		),
		nil,
		portBindings,
		&postgresUser,
	)
	if err != nil {
		return errors.New("CreateSupervisorConfError: " + err.Error())
	}

	hbaConfPath := "/etc/postgresql/" + versionStr + "/main/pg_hba.conf"
	_, err = infraHelper.RunCmdWithSubShell(
		"sed -i '1ilocal all all trust' " + hbaConfPath,
	)
	if err != nil {
		return errors.New("UpdatePgHbaError: " + err.Error())
	}

	err = SupervisordFacade{}.Reload()
	if err != nil {
		return errors.New("ReloadSupervisorError: " + err.Error())
	}

	time.Sleep(5 * time.Second)
	rootPass := infraHelper.GenPass(16)

	_, err = infraHelper.RunCmd(
		"psql",
		"-U",
		"postgres",
		"-c",
		"ALTER USER postgres WITH PASSWORD '"+rootPass+"';",
	)
	if err != nil {
		return errors.New("SetPostgresPassError: " + err.Error())
	}

	pgPassFileContent := "*:*:*:postgres:" + rootPass
	err = infraHelper.UpdateFile(
		"/root/.pgpass",
		pgPassFileContent,
		true,
	)
	if err != nil {
		return errors.New("CreatePgPassError: " + err.Error())
	}

	err = os.Chmod("/root/.pgpass", 0400)
	if err != nil {
		return errors.New("ChmodPgPassError: " + err.Error())
	}

	pgUserPgPassFilePath := "/var/lib/postgresql/.pgpass"
	err = infraHelper.UpdateFile(
		pgUserPgPassFilePath,
		pgPassFileContent,
		true,
	)
	if err != nil {
		return errors.New("CreatePgPassError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown",
		"postgres:postgres",
		pgUserPgPassFilePath,
	)
	if err != nil {
		return errors.New("ChownPgPassError: " + err.Error())
	}

	err = os.Chmod(pgUserPgPassFilePath, 0400)
	if err != nil {
		return errors.New("ChmodPgPassError: " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"sed -i '1s/.*/local all postgres scram-sha-256/' " + hbaConfPath,
	)
	if err != nil {
		return errors.New("UpdatePgHbaError: " + err.Error())
	}

	err = SupervisordFacade{}.Restart(createDto.Name)
	if err != nil {
		return errors.New("RestartPostgresqlError: " + err.Error())
	}

	return nil
}

func addRedis(createDto dto.CreateInstallableService) error {
	versionFlag := ""
	versionStr := "latest"
	if createDto.Version != nil {
		versionStr = createDto.Version.String()
		re := regexp.MustCompile(supportedServicesVersion["redis"])
		isVersionAllowed := re.MatchString(versionStr)

		if !isVersionAllowed {
			log.Printf("InvalidRedisVersion: %s", versionStr)
			return errors.New("InvalidRedisVersion")
		}
	}

	err := infraHelper.InstallPkgs(
		[]string{"lsb-release", "gpg"},
	)
	if err != nil {
		log.Printf("InstallPackagesError: %s", err)
		return errors.New("InstallPackagesError")
	}

	osRelease, err := infraHelper.GetOsRelease()
	if err != nil {
		log.Printf("GetOsReleaseError: %s", err)
		return errors.New("GetOsReleaseError")
	}

	err = installGpgKey("redis", "https://packages.redis.io/gpg")
	if err != nil {
		log.Printf("InstallGpgKeyError: %s", err)
		return errors.New("InstallGpgKeyError")
	}

	repoLine := "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb " + osRelease + " main"
	err = infraHelper.UpdateFile(
		"/etc/apt/sources.list.d/redis.list",
		repoLine,
		true,
	)
	if err != nil {
		log.Printf("CreateRepoFileError: %s", err)
		return errors.New("CreateRepoFileError")
	}

	if createDto.Version != nil {
		versionStr := createDto.Version.String()
		latestVersion, err := infraHelper.GetPkgLatestVersion(
			"redis-server",
			&versionStr,
		)
		if err != nil {
			log.Print(err)
			return err
		}

		versionFlag = "=" + latestVersion
	}

	err = infraHelper.InstallPkgs(
		[]string{RedisPackages[0] + versionFlag},
	)
	if err != nil {
		log.Printf("InstallServiceError: %s", err)
		return errors.New("InstallServiceError")
	}

	portBindings := []valueObject.PortBinding{
		valueObject.NewPortBinding(
			valueObject.NewNetworkPortPanic(6379),
			valueObject.NewNetworkProtocolPanic("tcp"),
		),
	}

	err = SupervisordFacade{}.CreateConf(
		createDto.Name,
		valueObject.NewServiceNaturePanic("solo"),
		valueObject.ServiceType("database"),
		valueObject.NewServiceVersionPanic(versionStr),
		valueObject.NewUnixCommandPanic("/usr/bin/redis-server /etc/redis/redis.conf"),
		nil,
		portBindings,
		nil,
	)
	if err != nil {
		return errors.New("CreateSupervisorConfError")
	}

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/^daemonize yes/daemonize no/g",
		"/etc/redis/redis.conf",
	)
	if err != nil {
		log.Printf("UpdateRedisConfError: %s", err)
		return errors.New("UpdateRedisConfError")
	}

	return nil
}

func CreateInstallable(
	createDto dto.CreateInstallableService,
) error {
	svcNameStr := createDto.Name.String()
	svcNameHasHash := strings.Contains(svcNameStr, "-")
	if svcNameHasHash {
		svcNameWithoutHash := strings.Split(svcNameStr, "-")[0]
		svcNameStr = svcNameWithoutHash
	}

	switch svcNameStr {
	case "php-webserver", "php":
		return addPhp()
	case "node":
		return addNode(createDto)
	case "mariadb":
		return addMariaDb(createDto)
	case "postgresql":
		return addPostgresqlDb(createDto)
	case "redis":
		return addRedis(createDto)
	default:
		return errors.New("UnknownInstallableService")
	}
}

func CreateInstallableSimplified(serviceName string) error {
	dto := dto.NewCreateInstallableService(
		valueObject.NewServiceNamePanic(serviceName),
		nil,
		nil,
		[]valueObject.PortBinding{},
		true,
	)
	return CreateInstallable(dto)
}
