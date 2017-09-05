package ts3

import (
	"net"
)

// Client is a TeamSpeak 3 ServerQuery client.
type Client struct {
	net.C
}

// Dial attempts to connect to a ServerQuery server.
func Dial(addr string) (*Client, error) {
	nc, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
}
