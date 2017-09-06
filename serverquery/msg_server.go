package serverquery

import (
	"context"
	"github.com/pkg/errors"
)

// UseCommand is the use command.
type UseCommand struct {
	// Port is the server to use.
	Port int `serverquery:"port"`
}

// GetResponseType returns an instance of the response type.
func (c *UseCommand) GetResponseType() interface{} {
	return nil
}

// GetCommandName returns the name of the command.
func (c *UseCommand) GetCommandName() string {
	return "use"
}

// UseServer selects a server by port
func (c *ServerQueryAPI) UseServer(ctx context.Context, port int) error {
	_, err := c.ExecuteCommand(ctx, &UseCommand{Port: port})
	return err
}

// LoginCommand is the login command.
type LoginCommand struct {
	username string
	password string
}

// GetResponseType returns an instance of the response type.
func (c *LoginCommand) GetResponseType() interface{} {
	return nil
}

// GetCommandName returns the name of the command.
func (c *LoginCommand) GetCommandName() string {
	return "login " + c.username + " " + c.password
}

// Login logs into the server.
func (c *ServerQueryAPI) Login(ctx context.Context, username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password must not be nil")
	}

	_, err := c.ExecuteCommand(ctx, &LoginCommand{username: username, password: password})
	return err
}
