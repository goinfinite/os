package terminalSessionInfra

import (
	"errors"
	"strings"
	"testing"

	tkInfra "github.com/goinfinite/tk/src/infra"
)

func TestParseSessionLine(t *testing.T) {
	testCases := []struct {
		name               string
		sessionLine        string
		expectedId         string
		expectedWorkingDir string
		expectedCommand    string
		expectedClients    uint16
		expectedErrMsg     string
	}{
		{
			name:               "ValidSessionLine",
			sessionLine:        "os-managed-0123456789abcdef|1700000000|2|/app|bash",
			expectedId:         "0123456789abcdef",
			expectedWorkingDir: "/app",
			expectedCommand:    "bash",
			expectedClients:    2,
		},
		{
			name:           "TooFewParts",
			sessionLine:    "os-managed-0123456789abcdef|1700000000|0|/app",
			expectedErrMsg: "InvalidSessionLine",
		},
		{
			name:           "TooManyParts",
			sessionLine:    "os-managed-0123456789abcdef|1700000000|0|/app|bash|extra",
			expectedErrMsg: "InvalidSessionLine",
		},
		{
			name:           "MissingSessionNamePrefix",
			sessionLine:    "other-0123456789abcdef|1700000000|0|/app|bash",
			expectedErrMsg: "NotATerminalSession",
		},
		{
			name:           "InvalidSessionId",
			sessionLine:    "os-managed-nothex|1700000000|0|/app|bash",
			expectedErrMsg: "InvalidTerminalSessionId",
		},
		{
			name:           "InvalidCreatedAt",
			sessionLine:    "os-managed-0123456789abcdef|notanumber|0|/app|bash",
			expectedErrMsg: "ParseSessionCreatedAtError",
		},
		{
			name:           "InvalidAttachedClients",
			sessionLine:    "os-managed-0123456789abcdef|1700000000|notanumber|/app|bash",
			expectedErrMsg: "ParseAttachedClientsError",
		},
		{
			name:           "EmptyWorkingDir",
			sessionLine:    "os-managed-0123456789abcdef|1700000000|0||bash",
			expectedErrMsg: "UnixAbsoluteFilePathValueMustNotBeEmpty",
		},
		{
			name:           "CommandTooShort",
			sessionLine:    "os-managed-0123456789abcdef|1700000000|0|/app|x",
			expectedErrMsg: "UnixCommandTooShort",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			client := &TerminalMultiplexerClient{}
			session, err := client.parseSessionLine(testCase.sessionLine)

			if testCase.expectedErrMsg != "" {
				if err == nil {
					t.Fatalf("Expected error containing '%s', got nil", testCase.expectedErrMsg)
				}
				if !strings.Contains(err.Error(), testCase.expectedErrMsg) {
					t.Fatalf(
						"Expected error containing '%s', got '%s'",
						testCase.expectedErrMsg, err.Error(),
					)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error, got '%s'", err.Error())
			}
			if session.Id.String() != testCase.expectedId {
				t.Errorf("Expected id '%s', got '%s'", testCase.expectedId, session.Id.String())
			}
			if session.WorkingDir.String() != testCase.expectedWorkingDir {
				t.Errorf(
					"Expected working dir '%s', got '%s'",
					testCase.expectedWorkingDir, session.WorkingDir.String(),
				)
			}
			if session.Command.String() != testCase.expectedCommand {
				t.Errorf(
					"Expected command '%s', got '%s'",
					testCase.expectedCommand, session.Command.String(),
				)
			}
			if session.AttachedClients != testCase.expectedClients {
				t.Errorf(
					"Expected %d attached clients, got %d",
					testCase.expectedClients, session.AttachedClients,
				)
			}
		})
	}
}

func TestIsServerAbsent(t *testing.T) {
	testCases := []struct {
		name           string
		multiplexerErr error
		expectedAbsent bool
	}{
		{
			name: "NoServerRunning",
			multiplexerErr: &tkInfra.ShellError{
				StdErr:   "no server running on /tmp/tmux-1000/default",
				ExitCode: 1,
			},
			expectedAbsent: true,
		},
		{
			name: "NoSessions",
			multiplexerErr: &tkInfra.ShellError{
				StdErr: "no sessions", ExitCode: 1,
			},
			expectedAbsent: true,
		},
		{
			name: "ErrorConnecting",
			multiplexerErr: &tkInfra.ShellError{
				StdErr:   "error connecting to /tmp/tmux-1000/default",
				ExitCode: 1,
			},
			expectedAbsent: true,
		},
		{
			name: "UnrelatedShellError",
			multiplexerErr: &tkInfra.ShellError{
				StdErr: "some other failure", ExitCode: 1,
			},
			expectedAbsent: false,
		},
		{
			name:           "NotAShellError",
			multiplexerErr: errors.New("SomeOtherError"),
			expectedAbsent: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			client := &TerminalMultiplexerClient{}
			isAbsent := client.isServerAbsent(testCase.multiplexerErr)
			if isAbsent != testCase.expectedAbsent {
				t.Errorf("Expected %t, got %t", testCase.expectedAbsent, isAbsent)
			}
		})
	}
}
