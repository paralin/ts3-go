package serverquery

import (
	"bytes"
)

// Command is an instance of a serverquery command.
type Command interface {
	// GetCommandName returns the command name.
	GetCommandName() string
	// GetResponseType returns an instance of the response type.
	GetResponseType() interface{}
}

// MarshalCommand marshals a command to a string.
func MarshalCommand(cmd Command, args ...interface{}) (string, error) {
	var result bytes.Buffer
	result.WriteString(cmd.GetCommandName())
	result.WriteRune(' ')
	a := []interface{}{cmd}
	a = append(a, args...)
	argStr, err := MarshalArguments(a...)
	if err != nil {
		return "", err
	}
	result.WriteString(argStr)
	return result.String(), nil
}
