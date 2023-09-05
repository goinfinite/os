package valueObject

import "errors"

type UnixCommand string

func NewUnixCommand(value string) (UnixCommand, error) {
	cmd := UnixCommand(value)
	if !cmd.isValid() {
		return "", errors.New("InvalidUnixCommand")
	}
	return cmd, nil
}

func NewUnixCommandPanic(value string) UnixCommand {
	cmd, err := NewUnixCommand(value)
	if err != nil {
		panic(err)
	}
	return cmd
}

func (cmd UnixCommand) isValid() bool {
	isTooShort := len(string(cmd)) < 3
	isTooLong := len(string(cmd)) > 2048
	return !isTooShort && !isTooLong
}

func (cmd UnixCommand) String() string {
	return string(cmd)
}
