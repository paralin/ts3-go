package serverquery

import (
	"bytes"
	"context"
	"strings"

	"github.com/pkg/errors"
)

// ServerQueryReadWriter can talk to the API.
type ServerQueryReadWriter interface {
	// WriteCommand writes a command line to the connection.
	WriteCommand(command string) error

	// ReadCommand reads a full command in.
	ReadCommand() (string, error)
}

// ServerQueryAPI is a client that implements the API.
type ServerQueryAPI struct {
	ServerQueryReadWriter

	// commandQueue is the queue of outgoing commands
	commandQueue chan *pendingCommand
	// readQueue contains incoming lines
	readQueue chan string
	// readError is the last read error
	readError error
}

// NewServerQueryAPI builds a new ServerQueryAPI client.
func NewServerQueryAPI(rw ServerQueryReadWriter) *ServerQueryAPI {
	return &ServerQueryAPI{
		ServerQueryReadWriter: rw,
		commandQueue:          make(chan *pendingCommand, 10),
		readQueue:             make(chan string, 10),
	}
}

// pendingCommand is a queued command waiting for a response.
type pendingCommand struct {
	ctx     context.Context
	command string
	result  interface{} // result is the result container
	doneCh  chan<- error
}

// ExecuteCommand sync-executes a command, waiting for the result.
func (a *ServerQueryAPI) ExecuteCommand(
	ctx context.Context,
	command string,
	result interface{},
) (interface{}, error) {
	doneCh := make(chan error, 1)
	pendCommand := &pendingCommand{
		ctx:     ctx,
		command: command,
		result:  result,
		doneCh:  doneCh,
	}
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case a.commandQueue <- pendCommand:
	}
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	case err, ok := <-doneCh:
		if !ok {
			return nil, context.Canceled
		}
		if err != nil {
			return nil, err
		}
		return pendCommand.result, nil
	}
}

// readPump reads messages from the connection.
func (a *ServerQueryAPI) readPump(ctx context.Context) (rerr error) {
	defer close(a.readQueue)
	defer func() { a.readError = rerr }()
	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		default:
		}
		cmd, err := a.ReadCommand()
		if err != nil {
			return context.Canceled
		}
		select {
		case a.readQueue <- cmd:
		case <-ctx.Done():
			return context.Canceled
		}
	}
}

// callResult is the result the server replies with
type callResult struct {
	// ErrorId is the error ID of the call
	ErrorId int `serverquery:"id"`
	// ErrorMessage is the message of the error.
	ErrorMessage string `serverquery:"msg"`
}

// submitCommand writes a command to the server and waits for the reply.
func (a *ServerQueryAPI) submitCommand(cmd string, resultObj interface{}) (interface{}, error) {
	if err := a.WriteCommand(cmd); err != nil {
		return nil, err
	}

	var resultBuf bytes.Buffer
	for {
		response, err := a.ReadCommand()
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(response, "error ") {
			response = response[len("error "):]
			respObj := &callResult{}

			ri, err := UnmarshalArguments(response[len("error "):], respObj)
			if err != nil {
				return nil, err
			}

			respObj = ri.(*callResult)
			if respObj.ErrorId != 0 {
				return nil, errors.Wrap(errors.New(respObj.ErrorMessage), "server error")
			}

			resultObj, err = UnmarshalArguments(resultBuf.String(), resultObj)
			return resultObj, nil
		}

		resultBuf.WriteString(" " + response)
	}
}

// processEvent processes and emits an event correctly.
func (a *ServerQueryAPI) processEvent(event string) {
	if !strings.HasPrefix(event, "notify") {
		return
	}
	return
}

// Run processes the client send/receive loop.
func (a *ServerQueryAPI) Run(parentContext context.Context) error {
	ctx, ctxCancel := context.WithCancel(parentContext)
	defer ctxCancel()
	go a.readPump(ctx)

	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case cmd := <-a.commandQueue:
			res, err := a.submitCommand(cmd.command, cmd.result)
			cmd.result = res
			select {
			case <-ctx.Done():
				return context.Canceled
			case cmd.doneCh <- err:
			}
		case env, ok := <-a.readQueue:
			if !ok {
				if a.readError == nil {
					return context.Canceled
				}
				return a.readError
			}
			a.processEvent(env)
		}
	}
}
