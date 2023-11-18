package servicesInfra

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"

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
		log.Printf("StopServiceError: %s", err)
		return errors.New("StopServiceError")
	}

	switch name.String() {
	case "php":
		infraHelper.RunCmd(
			"/usr/local/lsws/bin/lswsctrl",
			"stop",
		)
		infraHelper.RunCmd(
			"pkill",
			"lsphp",
		)
		infraHelper.RunCmd(
			"pkill",
			"sleep",
		)
	case "mysql":
		infraHelper.RunCmd(
			"mysqladmin",
			"shutdown",
		)
	case "node":
		infraHelper.RunCmd(
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
	switch name.String() {
	case "php":
		_, err := infraHelper.RunCmd(
			"/usr/local/lsws/bin/lswsctrl",
			"restart",
		)
		return err
	}

	_, err := infraHelper.RunCmd(
		supervisordCmd,
		"ctl",
		"restart",
		name.String(),
	)
	if err != nil {
		log.Printf("RestartServiceError: %s", err)
		return errors.New("RestartServiceError")
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
		log.Printf("ReloadSupervisorError: %s", err)
		return errors.New("ReloadSupervisorError")
	}

	return nil
}

func (facade SupervisordFacade) AddConf(svcName string, svcCmd string) error {
	svcConf := `
[program:` + svcName + `]
command=` + svcCmd + `
user=root
directory=/speedia
autostart=true
autorestart=true
startretries=3
startsecs=5
stdout_logfile=/dev/stdout
stdout_logfile_maxbytes=0
stderr_logfile=/dev/stderr
stderr_logfile_maxbytes=0
`

	f, err := os.OpenFile(supervisordConf, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("OpenSupervisorConfError: %s", err)
		return errors.New("OpenSupervisorConfError")
	}
	defer f.Close()

	if _, err := f.WriteString(svcConf); err != nil {
		log.Printf("WriteSupervisorConfError: %s", err)
		return errors.New("WriteSupervisorConfError")
	}

	return nil
}

func (facade SupervisordFacade) RemoveConf(svcName string) error {
	fileContent, err := os.ReadFile(supervisordConf)
	if err != nil {
		log.Printf("OpenSupervisorConfError: %s", err)
		return errors.New("OpenSupervisorConfError")
	}

	re := regexp.MustCompile(
		`\n?\[program:` + svcName + `\][\s\S]*?stderr_logfile_maxbytes=0\n?`,
	)
	updatedContent := re.ReplaceAll(fileContent, []byte{})

	err = os.WriteFile(supervisordConf, updatedContent, 0644)
	if err != nil {
		log.Printf("WriteSupervisorConfError: %s", err)
		return errors.New("WriteSupervisorConfError")
	}

	return nil
}
