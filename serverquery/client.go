package serverquery

import (
	"bytes"
	"context"
	"strings"
	"sync"

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
	// eventListeners is the list of event listeners
	eventListeners []chan<- Event
	// eventListenersMtx is the mtx of listeners
	eventListenersMtx sync.Mutex
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
	command Command,
) (interface{}, error) {
	doneCh := make(chan error, 1)
	mc, err := MarshalCommand(command)
	if err != nil {
		return nil, err
	}
	pendCommand := &pendingCommand{
		ctx:     ctx,
		command: mc,
		result:  command.GetResponseType(),
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

			if resultObj != nil {
				resultObj, err = UnmarshalArguments(resultBuf.String(), resultObj)
				if err != nil {
					return nil, err
				}
			}
			return resultObj, nil
		}

		resultBuf.WriteString(" " + response)
	}
}

// processEvent processes and emits an event correctly.
func (a *ServerQueryAPI) processEvent(ctx context.Context, event string) {
	msgParts := strings.SplitN(event, " ", 2)
	msgName := msgParts[0]
	if !strings.HasPrefix(msgName, "notify") {
		return
	}
	eventName := msgName[len("notify"):]
	c, ok := eventConstructorTable[eventName]
	if !ok {
		return
	}
	proto, err := UnmarshalArguments(msgParts[1], c())
	if err != nil {
		return
	}
	eve := proto.(Event)

	a.eventListenersMtx.Lock()
	for _, eventListener := range a.eventListeners {
		select {
		case <-ctx.Done():
			return
		case eventListener <- eve:
		}
	}
	a.eventListenersMtx.Unlock()
}

// Events returns a channel of events
func (a *ServerQueryAPI) Events() <-chan Event {
	ch := make(chan Event, 10)
	a.eventListenersMtx.Lock()
	a.eventListeners = append(a.eventListeners, ch)
	a.eventListenersMtx.Unlock()
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
