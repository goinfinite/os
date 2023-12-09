package servicesInfra

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
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
)

var supportedServicesVersion = map[string]string{
	"mariadb": `^(10\.([6-9]|10|11)|11\.[0-9]{1,2})$`,
	"node":    `^(1[2-9]|20)$`,
	"redis":   `^6\.(0|2)|^7\.0$`,
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

	virtualHost := os.Getenv("VIRTUAL_HOST")
	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/speedia.net/"+virtualHost+"/g",
		"/usr/local/lsws/conf/httpd_config.conf",
	)
	if err != nil {
		return errors.New("RenameHttpdVHostError: " + err.Error())
	}

	err = infraHelper.MakeDir("/app/conf/php")
	if err != nil {
		return errors.New("CreateConfDirError: " + err.Error())
	}

	err = copyAssets(
		"php/primary.conf",
		"/app/conf/php/primary.conf",
	)
	if err != nil {
		return errors.New("CopyAssetsError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/speedia.net/"+virtualHost+"/g",
		"/app/conf/php/primary.conf",
	)
	if err != nil {
		return errors.New("RenameVHostError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown",
		"-R",
		"lsadm:nogroup",
		"/app/conf/php",
	)
	if err != nil {
		return errors.New("ChownConfDirError: " + err.Error())
	}

	err = infraHelper.MakeDir("/app/logs/php")
	if err != nil {
		return errors.New("CreateLogDirError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown",
		"-R",
		"nobody:nogroup",
		"/app/logs/php",
	)
	if err != nil {
		return errors.New("ChownLogDirError: " + err.Error())
	}

	err = SupervisordFacade{}.AddConf(
		valueObject.NewServiceNamePanic("php"),
		valueObject.NewServiceNaturePanic("solo"),
		valueObject.NewServiceTypePanic("runtime"),
		valueObject.NewServiceVersionPanic("latest"),
		valueObject.NewUnixCommandPanic("/usr/local/lsws/bin/litespeed -d"),
		nil,
		[]valueObject.NetworkPort{
			valueObject.NewNetworkPortPanic(8080),
			valueObject.NewNetworkPortPanic(8443),
		},
	)
	if err != nil {
		return errors.New("AddSupervisorConfError: " + err.Error())
	}

	return nil
}

func addNode(addDto dto.AddInstallableService) error {
	versionStr := "lts"
	if addDto.Version != nil {
		versionStr = addDto.Version.String()
		re := regexp.MustCompile(supportedServicesVersion["node"])
		isVersionAllowed := re.MatchString(versionStr)

		if !isVersionAllowed {
			return errors.New("InvalidNodeVersion: " + versionStr)
		}
	}

	_, err := infraHelper.RunCmdWithSubShell(
		"rtx install node@" + versionStr,
	)
	if err != nil {
		return errors.New("InstallNodeError: " + err.Error())
	}

	appHtmlDir := "/app/html"
	err = infraHelper.MakeDir(appHtmlDir)
	if err != nil {
		return errors.New("CreateBaseDirError: " + err.Error())
	}

	startupFile := valueObject.NewUnixFilePathPanic(appHtmlDir + "/index.js")
	if addDto.StartupFile != nil {
		startupFile = *addDto.StartupFile
	}

	if !infraHelper.FileExists(startupFile.String()) {
		err = copyAssets(
			"nodejs/base-index.js",
			startupFile.String(),
		)
		if err != nil {
			return errors.New("CopyAssetsError: " + err.Error())
		}

		_, err = infraHelper.RunCmd(
			"chown",
			"nobody:nogroup",
			startupFile.String(),
		)
		if err != nil {
			return errors.New("ChownDummyIndexError: " + err.Error())
		}
	}

	ports := []valueObject.NetworkPort{
		valueObject.NewNetworkPortPanic(3000),
	}
	if len(addDto.Ports) > 0 {
		ports = addDto.Ports
	}

	startupFileBytes := []byte(startupFile.String())
	startupFileHash := md5.Sum(startupFileBytes)
	startupFileHashStr := hex.EncodeToString(startupFileHash[:])
	startupFileShortHashStr := startupFileHashStr[:12]

	svcNameWithSuffix := addDto.Name.String() + "-" + startupFileShortHashStr

	err = SupervisordFacade{}.AddConf(
		valueObject.NewServiceNamePanic(svcNameWithSuffix),
		valueObject.NewServiceNaturePanic("multi"),
		valueObject.NewServiceTypePanic("runtime"),
		valueObject.NewServiceVersionPanic(versionStr),
		valueObject.NewUnixCommandPanic(
			"rtx x node@"+versionStr+" -- node "+startupFile.String()+" &",
		),
		&startupFile,
		ports,
	)
	if err != nil {
		return errors.New("AddSupervisorConfError")
	}

	return nil
}

func addMariaDb(addDto dto.AddInstallableService) error {
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
	if addDto.Version != nil {
		versionStr = addDto.Version.String()
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
		"/usr/bin/mysqld_safe",
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
		"mysql",
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
		"mysqladmin",
		"shutdown",
	)
	if err != nil {
		return errors.New("StopMysqldSafeError: " + err.Error())
	}

	ports := []valueObject.NetworkPort{
		valueObject.NewNetworkPortPanic(3306),
	}
	if len(addDto.Ports) > 0 {
		ports = addDto.Ports
	}

	err = SupervisordFacade{}.AddConf(
		addDto.Name,
		valueObject.NewServiceNaturePanic("solo"),
		valueObject.NewServiceTypePanic("database"),
		valueObject.NewServiceVersionPanic(versionStr),
		valueObject.NewUnixCommandPanic("/usr/bin/mysqld_safe"),
		nil,
		ports,
	)
	if err != nil {
		return errors.New("AddSupervisorConfError")
	}

	return nil
}

func addRedis(addDto dto.AddInstallableService) error {
	versionFlag := ""
	versionStr := "latest"
	if addDto.Version != nil {
		versionStr = addDto.Version.String()
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

	err = infraHelper.DownloadFile(
		"https://packages.redis.io/gpg",
		"/speedia/redis.gpg",
	)
	if err != nil {
		log.Printf("DownloadRepoFileError: %s", err)
		return errors.New("DownloadRepoFileError")
	}

	_, err = infraHelper.RunCmd(
		"gpg",
		"--batch",
		"--yes",
		"--dearmor",
		"-o",
		"/usr/share/keyrings/redis-archive-keyring.gpg",
		"/speedia/redis.gpg",
	)
	if err != nil {
		log.Printf("GpgImportError: %s", err)
		return errors.New("GpgImportError")
	}

	err = os.Remove("/speedia/redis.gpg")
	if err != nil {
		log.Printf("RemoveRepoFileError: %s", err)
		return errors.New("RemoveRepoFileError")
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

	if addDto.Version != nil {
		versionStr := addDto.Version.String()
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

	err = SupervisordFacade{}.AddConf(
		addDto.Name,
		valueObject.NewServiceNaturePanic("solo"),
		valueObject.ServiceType("database"),
		valueObject.NewServiceVersionPanic(versionStr),
		valueObject.NewUnixCommandPanic("/usr/bin/redis-server /etc/redis/redis.conf"),
		nil,
		[]valueObject.NetworkPort{
			valueObject.NewNetworkPortPanic(6379),
		},
	)
	if err != nil {
		return errors.New("AddSupervisorConfError")
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

func AddInstallable(
	addDto dto.AddInstallableService,
) error {
	switch addDto.Name.String() {
	case "php":
		return addPhp()
	case "node":
		return addNode(addDto)
	case "mariadb":
		return addMariaDb(addDto)
	case "redis":
		return addRedis(addDto)
	default:
		return errors.New("UnknownInstallableService")
	}
}

func AddInstallableSimplified(serviceName string) error {
	dto := dto.NewAddInstallableService(
		valueObject.NewServiceNamePanic(serviceName),
		nil,
		nil,
		[]valueObject.NetworkPort{},
	)
	return AddInstallable(dto)
}
