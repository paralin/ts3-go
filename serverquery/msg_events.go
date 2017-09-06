package serverquery

import (
	"context"
)

// eventConstructorTable is the table of client event constructors.
var eventConstructorTable = make(map[string]func() Event)

// addEventPrototype adds an event constructor
func addEventPrototype(c func() Event) {
	p := c()
	eventConstructorTable[p.GetEventName()] = c
}

// Event is an instance of an event.
type Event interface {
	// GetEventName returns the event name.
	GetEventName() string
}

// UnknownEvent contains an unknown event.
type UnknownEvent struct {
	// EventSource is the contents of the event.
	EventSource string
}

// GetEventName returns the event name.
func (e *UnknownEvent) GetEventName() string {
	return "unknown"
}

// ClientLeftView is emitted when the client leaves the view.
type ClientLeftView struct {
	// SourceChannel is the channel the client left from.
	SourceChannel int `serverquery:"cfid"`
	// TargetChannel is the channel the client went to.
	TargetChannel int `serverquery:"ctid"`
	// ReasonId is the reason the client left.
	ReasonId int `serverquery:"reasonid"`
	// ReasonMessage is the reason why they are leaving.
	ReasonMessage string `serverquery:"reasonmsg"`
	// ClientId is the ID of the client.
	ClientId int `serverquery:"clid"`
}

// GetEventName returns the event name.
func (c *ClientLeftView) GetEventName() string {
	return "clientleftview"
}

func init() {
	addEventPrototype(func() Event {
		return &ClientLeftView{}
	})
}

// TextMessageReceived is emitted when the client receives a text message.
type TextMessageReceived struct {
	// TargetMode is the type of target.
	// 1 = user, 2 = channel 3 = server
	TargetMode int `serverquery:"targetmode"`
	// Message is the message received
	Message string `serverquery:"msg"`
	// TargetID is the ID of the target the message was sent to.
	TargetID int `serverquery:"target"`
	// InvokerID is the ID of the sender.
	InvokerID int `serverquery:"invokerid"`
	// InvokerName is the name of the sender.
	InvokerName string `serverquery:"invokername"`
	// InvokerUID is the unique identifier of the invoker.
	InvokerUID string `serverquery:"invokeruid"`
}

// GetEventName returns the event name.
func (c *TextMessageReceived) GetEventName() string {
	return "textmessage"
}

func init() {
	addEventPrototype(func() Event {
		return &TextMessageReceived{}
	})
}

// ServerNotifyRegisterCommand registers for events.
type ServerNotifyRegisterCommand struct {
	// EventType is one of server, channel, textserver, textchannel, textprivate
	EventType string `serverquery:"event"`
}

// GetResponseType returns an instance of the response type.
func (c *ServerNotifyRegisterCommand) GetResponseType() interface{} {
	return nil
}

// GetCommandName returns the name of the command.
func (c *ServerNotifyRegisterCommand) GetCommandName() string {
	return "servernotifyregister"
}

// ServerNotifyRegisterWithIdCommand registers for events with an ID.
type ServerNotifyRegisterWithIdCommand struct {
	ServerNotifyRegisterCommand

	// Id limits to a specific ID.
	Id int `serverquery:"id"`
}

// GetResponseType returns an instance of the response type.
func (c *ServerNotifyRegisterWithIdCommand) GetResponseType() interface{} {
	return nil
}

// GetCommandName returns the name of the command.
func (c *ServerNotifyRegisterWithIdCommand) GetCommandName() string {
	return "servernotifyregister"
}

// ServerNotifyRegister registers for events, optionally with an id.
func (c *ServerQueryAPI) ServerNotifyRegister(
	ctx context.Context,
	eventType string,
	id int,
) error {
	var cmd Command
	snr := &ServerNotifyRegisterCommand{EventType: eventType}
	if id != 0 {
		cmd = &ServerNotifyRegisterWithIdCommand{
			ServerNotifyRegisterCommand: *snr,
			Id: id,
		}
	} else {
		cmd = snr
	}

	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// ServerNotifyRegisterAll registers for all ambient events
func (c *ServerQueryAPI) ServerNotifyRegisterAll(ctx context.Context) error {
	if err := c.ServerNotifyRegister(ctx, "server", 0); err != nil {
		return err
	}
	if err := c.ServerNotifyRegister(ctx, "channel", 0); err != nil {
		return err
	}
	if err := c.ServerNotifyRegister(ctx, "textserver", 0); err != nil {
		return err
	}
	if err := c.ServerNotifyRegister(ctx, "textchannel", 0); err != nil {
		return err
	}
	if err := c.ServerNotifyRegister(ctx, "textprivate", 0); err != nil {
		return err
	}
	return nil
}
