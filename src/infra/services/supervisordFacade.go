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
	name valueObject.ServiceName,
	isAutoStart bool,
) error {
	autoStartStr := strconv.FormatBool(isAutoStart)

	_, err := infraHelper.RunCmd(
		"sed",
		"-i",
		"-e",
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
	_, err := infraHelper.RunCmd(
		supervisordCmd,
		"ctl",
		"start",
		name.String(),
	)
	if err != nil {
		return errors.New("StartServiceError: " + err.Error())
	}

	err = facade.toggleAutoStart(name, true)
	if err != nil {
		return err
	}

	return nil
}

func (facade SupervisordFacade) Stop(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd(
		supervisordCmd,
		"ctl",
		"stop",
		name.String(),
	)
	if err != nil {
		return errors.New("StopServiceError: " + err.Error())
	}

	switch name.String() {
	case "nginx":
		_, _ = infraHelper.RunCmd(
			"pkill",
			"nginx",
		)
	case "php":
		_, _ = infraHelper.RunCmd(
			"/usr/local/lsws/bin/lswsctrl",
			"stop",
		)
		_, _ = infraHelper.RunCmd(
			"pkill",
			"lsphp",
		)
		_, _ = infraHelper.RunCmd(
			"pkill",
			"sleep",
		)
	case "mysql":
		_, _ = infraHelper.RunCmd(
			"mysqladmin",
			"shutdown",
		)
	case "node":
		_, _ = infraHelper.RunCmd(
			"pkill",
			"node",
		)
	}

	err = facade.toggleAutoStart(name, false)
	if err != nil {
		return err
	}

	return nil
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
	_, err := infraHelper.RunCmd(
		supervisordCmd,
		"ctl",
		"reload",
	)
	if err != nil {
		return errors.New("ReloadSupervisorError: " + err.Error())
	}

	return nil
}

func (facade SupervisordFacade) AddConf(
	svcName string,
	svcCmd string,
	svcType string,
	svcPorts []int,
) error {
	err := infraHelper.MakeDir("/app/logs/" + svcName)
	if err != nil {
		return errors.New("CreateLogDirError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(
		"chown",
		"-R",
		"nobody:nogroup",
		"/app/logs/"+svcName,
	)
	if err != nil {
		return errors.New("ChownLogDirError: " + err.Error())
	}

	svcType = "SVC_TYPE=\"" + svcType + "\""

	svcPortsStr := ""
	if len(svcPorts) > 0 {
		portsStrSlice := []string{}
		for _, port := range svcPorts {
			portsStrSlice = append(portsStrSlice, strconv.Itoa(port))
		}
		svcPortsStr = ",SVC_PORTS=\"" + strings.Join(portsStrSlice, ",") + "\""
	}

	logFilePath := "/app/logs/" + svcName + "/" + svcName + ".log"

	svcConf := `
[program:` + svcName + `]
command=` + svcCmd + `
user=root
directory=/speedia
autostart=true
autorestart=true
startretries=3
startsecs=3
stdout_logfile=` + logFilePath + `
stdout_logfile_maxbytes=10MB
stderr_logfile=` + logFilePath + `
stderr_logfile_maxbytes=10MB
environment=` + svcType + svcPortsStr + `
`

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

func (facade SupervisordFacade) RemoveConf(svcName string) error {
	fileContent, err := os.ReadFile(supervisordConf)
	if err != nil {
		return errors.New("OpenSupervisorConfError: " + err.Error())
	}

	re := regexp.MustCompile(
		`\n?\[program:` + svcName + `\][\s\S]*?stderr_logfile_maxbytes=0\n?`,
	)
	updatedContent := re.ReplaceAll(fileContent, []byte{})

	err = os.WriteFile(supervisordConf, updatedContent, 0644)
	if err != nil {
		return errors.New("WriteSupervisorConfError: " + err.Error())
	}

	return nil
}
