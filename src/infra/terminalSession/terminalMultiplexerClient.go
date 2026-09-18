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

	terminalMultiplexerSessionListFormat = "#{session_name}|#{session_created}|" +
		"#{session_attached}|#{pane_current_path}|#{pane_current_command}"
)

var terminalMultiplexerServerAbsentRegex = regexp.MustCompile(
	`(?im)^(no server running|no sessions|error connecting to).*$`,
)

type TerminalMultiplexerSession struct {
	Id              valueObject.TerminalSessionId
	CreatedAt       tkValueObject.UnixTime
	AttachedClients uint16
	WorkingDir      tkValueObject.UnixAbsoluteFilePath
	Command         tkValueObject.UnixCommand
}

type TerminalMultiplexerClient struct {
	accountUsername valueObject.Username
	userId          uint32
	groupId         uint32
}

func NewTerminalMultiplexerClient(
	accountUsername valueObject.Username,
) (*TerminalMultiplexerClient, error) {
	userInfo, err := user.Lookup(accountUsername.String())
	if err != nil {
		return nil, errors.New("AccountUserLookupError: " + err.Error())
	}

	userId, err := strconv.ParseUint(userInfo.Uid, 10, 32)
	if err != nil {
		return nil, errors.New("ParseUserIdError: " + err.Error())
	}

	groupId, err := strconv.ParseUint(userInfo.Gid, 10, 32)
	if err != nil {
		return nil, errors.New("ParseGroupIdError: " + err.Error())
	}

	return &TerminalMultiplexerClient{
		accountUsername: accountUsername,
		userId:          uint32(userId),
		groupId:         uint32(groupId),
	}, nil
}

func (client *TerminalMultiplexerClient) buildSessionName(
	terminalSessionId valueObject.TerminalSessionId,
) string {
	return TerminalSessionNamePrefix + terminalSessionId.String()
}

func (client *TerminalMultiplexerClient) accountShellEnvironment() []string {
	return []string{
		"HOME=" + infraEnvs.ApplicationRootDir,
		"USER=" + client.accountUsername.String(),
		"LOGNAME=" + client.accountUsername.String(),
		"SHELL=/bin/bash",
		"TERM=xterm-256color",
	}
}

func (client *TerminalMultiplexerClient) runMultiplexerCommand(
	command string,
	args []string,
) (string, error) {
	shellSettings := tkInfra.ShellSettings{
		Command:           command,
		Args:              args,
		ShouldUseCleanEnv: true,
		Username:          client.accountUsername.String(),
		Envs:              client.accountShellEnvironment(),
	}

	return tkInfra.NewShell(shellSettings).Run()
}

func (client *TerminalMultiplexerClient) CreateSession(
	terminalSessionId valueObject.TerminalSessionId,
	workingDir tkValueObject.UnixAbsoluteFilePath,
) error {
	_, err := client.runMultiplexerCommand("tmux", []string{
		"new-session", "-d", "-s", client.buildSessionName(terminalSessionId),
		"-c", workingDir.String(),
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
	// Example session line: os-managed-0123456789abcdef|1700000000|2|/app|bash
	lineParts := strings.Split(sessionLine, "|")
	if len(lineParts) != 5 {
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

	return TerminalMultiplexerSession{
		Id:              terminalSessionId,
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
	ptyFile   *os.File
	attachCmd *exec.Cmd
	waitOnce  sync.Once
	waitErr   error
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

func (handle *terminalMultiplexerAttachHandle) Close() error {
	closeErr := handle.ptyFile.Close()

	killErr := handle.attachCmd.Process.Kill()
	if killErr != nil {
		slog.Debug("KillAttachProcessError", slog.String("err", killErr.Error()))
	}

	waitErr := handle.Wait()
	if waitErr != nil {
		slog.Debug("WaitAttachProcessError", slog.String("err", waitErr.Error()))
	}

	return closeErr
}

func (client *TerminalMultiplexerClient) Attach(
	terminalSessionId valueObject.TerminalSessionId,
) (repository.TerminalSessionAttachHandle, error) {
	attachCmd := exec.Command(
		"tmux", "attach", "-t", client.buildSessionName(terminalSessionId),
	)
	attachCmd.Dir = infraEnvs.ApplicationRootDir
	attachCmd.Env = append(
		client.accountShellEnvironment(),
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
		ptyFile:   ptyFile,
		attachCmd: attachCmd,
	}, nil
}
