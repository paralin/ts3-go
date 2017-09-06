package serverquery

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
