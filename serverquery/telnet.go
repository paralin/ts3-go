package serverquery

import (
	"bytes"
	"io"
	"net"
)

// ServerQueryReadWriter can talk to the API
type ServerQueryReadWriter struct {
	net.Conn
}

// WriteCommand writes a command line to the connection.
func (rw *ServerQueryReadWriter) WriteCommand(command string) error {
	var buf bytes.Buffer
	buf.WriteString(command)
	buf.WriteRune('\n')
	_, err := io.Copy(rw.Conn, &buf)
	return err
}

// ReadCommand reads a full command in.
func (rw *ServerQueryReadWriter) ReadCommand() (string, error) {
	var result bytes.Buffer
	buf := make([]byte, 1)
	for {
		_, err := rw.Conn.Read(buf)
		if err != nil {
			return "", err
		}
		r := rune(buf[0])
		if r == '\n' {
			return result.String(), nil
		}
		result.WriteRune(r)
	}
}
