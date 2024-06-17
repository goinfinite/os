package servicesInfra

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type SupervisordFacade struct {
}

const supervisordCmd string = "/usr/bin/supervisord"
const supervisordConf string = "/speedia/supervisord.conf"

func (facade SupervisordFacade) toggleAutoStart(
	name valueObject.ServiceName, shouldAutoStart bool,
) error {
	autoStartStr := strconv.FormatBool(shouldAutoStart)

	_, err := infraHelper.RunCmd(
		"sed", "-i", "-e",
		"/\\[program:"+name.String()+"\\]/,"+
			"/^\\[/{s/autostart=.*/autostart="+autoStartStr+"/"+
			";s/autorestart=.*/autorestart="+autoStartStr+"/}",
		supervisordConf,
	)
	if err != nil {
		return errors.New("UpdateSupervisordConfError: " + err.Error())
	}

	return nil
}

func (facade SupervisordFacade) Start(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd(supervisordCmd, "ctl", "start", name.String())
	if err != nil {
		return errors.New("StartServiceError: " + err.Error())
	}

	err = facade.toggleAutoStart(name, true)
	if err != nil {
		return err
	}

	return nil
}

func (facade SupervisordFacade) stopServiceByName(svcName string) error {
	_, err := infraHelper.RunCmd(supervisordCmd, "ctl", "stop", svcName)
	if err != nil {
		return errors.New("StopServiceError: " + err.Error())
	}

	return nil
}

func (facade SupervisordFacade) stopNginx() error {
	_, err := infraHelper.RunCmdWithSubShell("nginx -t")
	if err != nil {
		return errors.New("NginxTestFailed: " + err.Error())
	}

	err = facade.stopServiceByName("nginx")
	if err != nil {
		return err
	}

	_, _ = infraHelper.RunCmd("pkill", "nginx")
	return nil
}

func (facade SupervisordFacade) stopPhpWebServer() error {
	err := facade.stopServiceByName("php-webserver")
	if err != nil {
		return err
	}

	_, _ = infraHelper.RunCmd("/usr/local/lsws/bin/lswsctrl", "stop")
	_, _ = infraHelper.RunCmd("pkill", "lsphp")
	_, _ = infraHelper.RunCmd("pkill", "sleep")
	return nil
}

func (facade SupervisordFacade) stopMariaDb() error {
	err := facade.stopServiceByName("mariadb")
	if err != nil {
		return err
	}

	_, _ = infraHelper.RunCmd("mysqladmin", "--defaults-file=/root/.my.cnf", "shutdown")
	return nil
}

func (facade SupervisordFacade) Stop(name valueObject.ServiceName) error {
	switch name.String() {
	case "nginx":
		err := facade.stopNginx()
		if err != nil {
			return err
		}
	case "php-webserver", "php":
		err := facade.stopPhpWebServer()
		if err != nil {
			return err
		}
	case "mariadb", "mysql":
		err := facade.stopMariaDb()
		if err != nil {
			return err
		}
	}

	return facade.toggleAutoStart(name, false)
}

func (facade SupervisordFacade) Restart(name valueObject.ServiceName) error {
	err := facade.Stop(name)
	if err != nil {
		return errors.New("StopServiceError: " + err.Error())
	}

	time.Sleep(3 * time.Second)

	err = facade.Start(name)
	if err != nil {
		return errors.New("StartServiceError: " + err.Error())
	}

	return nil
}

func (facade SupervisordFacade) Reload() error {
	_, err := infraHelper.RunCmd(supervisordCmd, "ctl", "reload")
	if err != nil {
		return errors.New("ReloadSupervisorError: " + err.Error())
	}

	return nil
}

func (facade SupervisordFacade) CreateConf(
	svcName valueObject.ServiceName,
	svcNature valueObject.ServiceNature,
	svcType valueObject.ServiceType,
	svcVersion valueObject.ServiceVersion,
	svcCmd valueObject.UnixCommand,
	startupFile *valueObject.UnixFilePath,
	svcPortBindings []valueObject.PortBinding,
	svcUser *valueObject.Username,
) error {
	svcNameStr := svcName.String()

	err := infraHelper.MakeDir("/app/logs/" + svcNameStr)
	if err != nil {
		return errors.New("CreateLogDirError: " + err.Error())
	}

	chownRecursively := true
	chownSymlinksToo := false
	err = infraHelper.UpdatePermissionsForWebServerUse(
		"/app/logs/"+svcNameStr, chownRecursively, chownSymlinksToo,
	)
	if err != nil {
		return errors.New("ChownLogDirError: " + err.Error())
	}

	svcNatureStr := "SVC_NATURE=\"" + svcNature.String() + "\""

	svcTypeStr := ",SVC_TYPE=\"" + svcType.String() + "\""

	svcVersionStr := ",SVC_VERSION=\"" + svcVersion.String() + "\""

	startupFileStr := ""
	if startupFile != nil {
		startupFileStr = ",SVC_STARTUP_FILE=\"" + startupFile.String() + "\""
	}

	svcPortBindingsStr := ""
	if len(svcPortBindings) > 0 {
		portBindingsStrSlice := []string{}
		for _, portBinding := range svcPortBindings {
			portBindingsStrSlice = append(
				portBindingsStrSlice,
				portBinding.String(),
			)
		}
		svcPortBindingsStr = ",SVC_PORT_BINDINGS=\"" +
			strings.Join(portBindingsStrSlice, ",") + "\""
	}

	svcUserStr := "root"
	if svcUser != nil {
		svcUserStr = svcUser.String()
	}

	logFilePath := "/app/logs/" + svcNameStr + "/" + svcNameStr + ".log"

	// cSpell:disable
	svcConf := `
[program:` + svcNameStr + `]
command=` + svcCmd.String() + `
user=` + svcUserStr + `
directory=/speedia
autostart=true
autorestart=true
startretries=3
startsecs=3
stdout_logfile=` + logFilePath + `
stdout_logfile_maxbytes=10MB
stderr_logfile=` + logFilePath + `
stderr_logfile_maxbytes=10MB
environment=` + svcNatureStr + svcTypeStr + svcVersionStr + startupFileStr + svcPortBindingsStr + `
`
	// cSpell:enable

	f, err := os.OpenFile(supervisordConf, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("OpenSupervisorConfError: " + err.Error())
	}
	defer f.Close()

	if _, err := f.WriteString(svcConf); err != nil {
		return errors.New("WriteSupervisorConfError: " + err.Error())
	}

	return nil
}

func (facade SupervisordFacade) RemoveConf(svcName valueObject.ServiceName) error {
	fileContent, err := os.ReadFile(supervisordConf)
	if err != nil {
		return errors.New("OpenSupervisorConfError: " + err.Error())
	}

	re := regexp.MustCompile(
		`\n?\[program:` + svcName.String() + `\][\s\S]*?environment=[^\n]*\n?`,
	)
	updatedContent := re.ReplaceAll(fileContent, []byte{})

	err = os.WriteFile(supervisordConf, updatedContent, 0644)
	if err != nil {
		return errors.New("WriteSupervisorConfError: " + err.Error())
	}

	return nil
}
