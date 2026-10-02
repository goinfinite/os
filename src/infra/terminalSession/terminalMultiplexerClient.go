package terminalSessionInfra

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/creack/pty"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

const (
	TerminalSessionNamePrefix = "os-managed-"

	terminalMultiplexerNameOption        = "@os-name"
	terminalMultiplexerSocketNamePrefix  = "os-"
	terminalMultiplexerSessionListFormat = "#{session_name}|#{session_created}|" +
		"#{session_attached}|#{pane_current_path}|#{pane_current_command}|" +
		"#{" + terminalMultiplexerNameOption + "}"
)

var terminalMultiplexerServerAbsentRegex = regexp.MustCompile(
	`(?im)^(no server running|no sessions|error connecting to).*$`,
)

type TerminalMultiplexerSession struct {
	Id              valueObject.TerminalSessionId
	Name            *valueObject.TerminalSessionName
	CreatedAt       tkValueObject.UnixTime
	AttachedClients uint16
	WorkingDir      tkValueObject.UnixAbsoluteFilePath
	Command         tkValueObject.UnixCommand
}

type TerminalMultiplexerClient struct {
	accountUsername valueObject.Username
	runAsUsername   valueObject.Username
	tmuxSocketName  string
	userId          uint32
	groupId         uint32
}

func NewTerminalMultiplexerClient(
	accountUsername, runAsUsername valueObject.Username,
) (*TerminalMultiplexerClient, error) {
	userInfo, err := user.Lookup(runAsUsername.String())
	if err != nil {
		return nil, errors.New("RunAsUserLookupError: " + err.Error())
	}

	userId, err := strconv.ParseUint(userInfo.Uid, 10, 32)
	if err != nil {
		return nil, errors.New("ParseUserIdError: " + err.Error())
	}

	groupId, err := strconv.ParseUint(userInfo.Gid, 10, 32)
	if err != nil {
		return nil, errors.New("ParseGroupIdError: " + err.Error())
	}

	tmuxSocketName := ""
	if runAsUsername != accountUsername {
		tmuxSocketName = terminalMultiplexerSocketNamePrefix +
			accountUsername.String() + "-" + runAsUsername.String()
	}

	return &TerminalMultiplexerClient{
		accountUsername: accountUsername,
		runAsUsername:   runAsUsername,
		tmuxSocketName:  tmuxSocketName,
		userId:          uint32(userId),
		groupId:         uint32(groupId),
	}, nil
}

func (client *TerminalMultiplexerClient) buildSessionName(
	terminalSessionId valueObject.TerminalSessionId,
) string {
	return TerminalSessionNamePrefix + terminalSessionId.String()
}

func (client *TerminalMultiplexerClient) buildAccountShellEnvironment() []string {
	runAsUsernameStr := client.runAsUsername.String()
	environment := []string{
		"USER=" + runAsUsernameStr,
		"LOGNAME=" + runAsUsernameStr,
		"SHELL=/bin/bash",
		"TERM=xterm-256color",
	}

	accountHomeDir := infraEnvs.UserDataBaseDirectory + "/" + client.accountUsername.String()
	if client.runAsUsername == valueObject.UsernameNobody {
		return append(environment,
			"HOME="+infraEnvs.ApplicationRootDir,
			"MISE_DATA_DIR="+infraEnvs.ToolchainDataDir,
			"MISE_CONFIG_DIR="+infraEnvs.ToolchainDataDir+"/config",
			"MISE_STATE_DIR="+infraEnvs.ToolchainDataDir+"/state",
			"MISE_CACHE_DIR="+infraEnvs.ToolchainDataDir+"/cache",
		)
	}

	return append(environment,
		"HOME="+accountHomeDir,
		"MISE_DATA_DIR="+accountHomeDir+"/.local/share/mise",
		"MISE_CONFIG_DIR="+accountHomeDir+"/.config/mise",
		"MISE_STATE_DIR="+accountHomeDir+"/.local/state/mise",
		"MISE_CACHE_DIR="+accountHomeDir+"/.cache/mise",
	)
}

func (client *TerminalMultiplexerClient) buildMultiplexerSocketArgs(
	args []string,
) []string {
	if client.tmuxSocketName == "" {
		return args
	}

	return append([]string{"-L", client.tmuxSocketName}, args...)
}

func (client *TerminalMultiplexerClient) runMultiplexerCommand(
	command string,
	args []string,
) (string, error) {
	shellSettings := tkInfra.ShellSettings{
		Command:           command,
		Args:              client.buildMultiplexerSocketArgs(args),
		ShouldUseCleanEnv: true,
		Username:          client.runAsUsername.String(),
		Envs:              client.buildAccountShellEnvironment(),
	}

	return tkInfra.NewShell(shellSettings).Run()
}

func (client *TerminalMultiplexerClient) CreateSession(
	terminalSessionId valueObject.TerminalSessionId,
	workingDir tkValueObject.UnixAbsoluteFilePath,
) error {
	_, err := client.runMultiplexerCommand("tmux", []string{
		"new-session", "-d", "-s", client.buildSessionName(terminalSessionId),
		"-c", workingDir.String(), "-e", "COLORTERM=truecolor",
	})
	if err != nil {
		return errors.New("CreateSessionError: " + err.Error())
	}

	return nil
}

func (client *TerminalMultiplexerClient) SendCommand(
	terminalSessionId valueObject.TerminalSessionId,
	command tkValueObject.UnixCommand,
) error {
	_, err := client.runMultiplexerCommand("tmux", []string{
		"send-keys", "-t", client.buildSessionName(terminalSessionId),
		command.String(), "Enter",
	})
	if err != nil {
		return errors.New("SendCommandError: " + err.Error())
	}

	return nil
}

func (client *TerminalMultiplexerClient) SetSessionName(
	terminalSessionId valueObject.TerminalSessionId,
	name valueObject.TerminalSessionName,
) error {
	_, err := client.runMultiplexerCommand("tmux", []string{
		"set-option", "-t", client.buildSessionName(terminalSessionId),
		terminalMultiplexerNameOption, name.String(),
	})
	if err != nil {
		return errors.New("SetSessionNameError: " + err.Error())
	}

	return nil
}

func (client *TerminalMultiplexerClient) ClearSessionName(
	terminalSessionId valueObject.TerminalSessionId,
) error {
	_, err := client.runMultiplexerCommand("tmux", []string{
		"set-option", "-t", client.buildSessionName(terminalSessionId),
		"-u", terminalMultiplexerNameOption,
	})
	if err != nil {
		return errors.New("ClearSessionNameError: " + err.Error())
	}

	return nil
}

func (client *TerminalMultiplexerClient) KillSession(
	terminalSessionId valueObject.TerminalSessionId,
) error {
	_, err := client.runMultiplexerCommand("tmux", []string{
		"kill-session", "-t", client.buildSessionName(terminalSessionId),
	})
	if err != nil {
		return errors.New("KillSessionError: " + err.Error())
	}

	return nil
}

func (client *TerminalMultiplexerClient) isServerAbsent(err error) bool {
	shellError, isShellError := err.(*tkInfra.ShellError)
	if !isShellError {
		return false
	}

	return terminalMultiplexerServerAbsentRegex.MatchString(shellError.StdErr)
}

func (client *TerminalMultiplexerClient) parseSessionLine(
	sessionLine string,
) (session TerminalMultiplexerSession, err error) {
	// Example session line: os-managed-0123456789abcdef|1700000000|2|/app|bash|opencode
	lineParts := strings.SplitN(sessionLine, "|", 6)
	if len(lineParts) != 6 {
		return session, errors.New("InvalidSessionLine")
	}

	sessionName := lineParts[0]
	if !strings.HasPrefix(sessionName, TerminalSessionNamePrefix) {
		return session, errors.New("NotATerminalSession")
	}

	terminalSessionId, err := valueObject.NewTerminalSessionId(
		strings.TrimPrefix(sessionName, TerminalSessionNamePrefix),
	)
	if err != nil {
		return session, err
	}

	createdAtUnixSecs, err := strconv.ParseInt(lineParts[1], 10, 64)
	if err != nil {
		return session, errors.New("ParseSessionCreatedAtError")
	}

	createdAt, err := tkValueObject.NewUnixTime(createdAtUnixSecs)
	if err != nil {
		return session, err
	}

	attachedClients, err := strconv.ParseUint(lineParts[2], 10, 16)
	if err != nil {
		return session, errors.New("ParseAttachedClientsError")
	}

	workingDir, err := tkValueObject.NewUnixAbsoluteFilePath(lineParts[3], true)
	if err != nil {
		return session, err
	}

	command, err := tkValueObject.NewUnixCommand(lineParts[4])
	if err != nil {
		return session, err
	}

	var namePtr *valueObject.TerminalSessionName
	if sessionName := strings.TrimSpace(lineParts[5]); sessionName != "" {
		sessionNameVo, err := valueObject.NewTerminalSessionName(sessionName)
		if err != nil {
			return session, err
		}
		namePtr = &sessionNameVo
	}

	return TerminalMultiplexerSession{
		Id:              terminalSessionId,
		Name:            namePtr,
		CreatedAt:       createdAt,
		AttachedClients: uint16(attachedClients),
		WorkingDir:      workingDir,
		Command:         command,
	}, nil
}

func (client *TerminalMultiplexerClient) ListSessions() (
	sessions []TerminalMultiplexerSession, err error,
) {
	stdoutStr, err := client.runMultiplexerCommand("tmux", []string{
		"list-sessions", "-F", terminalMultiplexerSessionListFormat,
	})
	if err != nil {
		if client.isServerAbsent(err) {
			return []TerminalMultiplexerSession{}, nil
		}

		return sessions, errors.New("ListSessionsError: " + err.Error())
	}

	sessions = []TerminalMultiplexerSession{}
	for sessionLine := range strings.SplitSeq(stdoutStr, "\n") {
		sessionLine = strings.TrimSpace(sessionLine)
		if sessionLine == "" {
			continue
		}

		session, err := client.parseSessionLine(sessionLine)
		if err != nil {
			slog.Debug(
				"ParseSessionLineError",
				slog.String("sessionLine", sessionLine),
				slog.String("err", err.Error()),
			)
			continue
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

type terminalMultiplexerAttachHandle struct {
	attachCmd *exec.Cmd
	ptyFile   *os.File
	sessionId valueObject.TerminalSessionId

	closeErr  error
	closeOnce sync.Once
	waitErr   error
	waitOnce  sync.Once
}

func (handle *terminalMultiplexerAttachHandle) Read(buffer []byte) (int, error) {
	return handle.ptyFile.Read(buffer)
}

func (handle *terminalMultiplexerAttachHandle) Write(buffer []byte) (int, error) {
	return handle.ptyFile.Write(buffer)
}

func (handle *terminalMultiplexerAttachHandle) Resize(cols, rows uint16) error {
	return pty.Setsize(handle.ptyFile, &pty.Winsize{Cols: cols, Rows: rows})
}

func (handle *terminalMultiplexerAttachHandle) Wait() error {
	handle.waitOnce.Do(func() {
		handle.waitErr = handle.attachCmd.Wait()
	})

	return handle.waitErr
}

func (handle *terminalMultiplexerAttachHandle) hasKillSignalExit() bool {
	processState := handle.attachCmd.ProcessState
	if processState == nil {
		return false
	}

	waitStatus, isWaitStatus := processState.Sys().(syscall.WaitStatus)
	if !isWaitStatus {
		return false
	}

	return waitStatus.Signaled() && waitStatus.Signal() == syscall.SIGKILL
}

func (handle *terminalMultiplexerAttachHandle) logAttachProcessEnd(
	killErr, waitErr error,
) {
	if handle.hasKillSignalExit() {
		return
	}

	exitStatus := "exit status 0"
	if waitErr != nil {
		exitStatus = waitErr.Error()
	}

	logContext := []any{
		slog.String("sessionId", handle.sessionId.String()),
		slog.Int("pid", handle.attachCmd.Process.Pid),
		slog.String("exitStatus", exitStatus),
	}
	if killErr != nil {
		logContext = append(logContext, slog.String("killErr", killErr.Error()))
	}

	slog.Debug("AttachProcessEndedWithoutKillSignal", logContext...)
}

func (handle *terminalMultiplexerAttachHandle) Close() error {
	handle.closeOnce.Do(func() {
		killErr := handle.attachCmd.Process.Kill()
		handle.closeErr = handle.ptyFile.Close()
		waitErr := handle.Wait()
		handle.logAttachProcessEnd(killErr, waitErr)
	})

	return handle.closeErr
}

func (client *TerminalMultiplexerClient) BuildAttachArgs(
	terminalSessionId valueObject.TerminalSessionId,
) []string {
	return client.buildMultiplexerSocketArgs([]string{
		"attach", "-t", client.buildSessionName(terminalSessionId),
	})
}

func (client *TerminalMultiplexerClient) Attach(
	terminalSessionId valueObject.TerminalSessionId,
) (repository.TerminalSessionAttachHandle, error) {
	attachArgs := client.BuildAttachArgs(terminalSessionId)
	attachCmd := exec.Command("tmux", attachArgs...)
	attachCmd.Dir = infraEnvs.ApplicationRootDir
	attachCmd.Env = append(
		client.buildAccountShellEnvironment(),
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
	)
	attachCmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: client.userId,
			Gid: client.groupId,
		},
	}

	ptyFile, err := pty.Start(attachCmd)
	if err != nil {
		return nil, errors.New("StartAttachPtyError: " + err.Error())
	}

	return &terminalMultiplexerAttachHandle{
		attachCmd: attachCmd,
		ptyFile:   ptyFile,
		sessionId: terminalSessionId,
	}, nil
}
