package serverquery

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"regexp"
)

var invalidCharRegex = regexp.MustCompile(`[^ -~]+`)

// ServerQueryReadWriter can talk to the API
type ServerQueryReadWriter struct {
	net.Conn
	scanner *bufio.Scanner
}

// NewServerQueryReadWriter builds a new read-writer
func NewServerQueryReadWriter(conn net.Conn) *ServerQueryReadWriter {
	return &ServerQueryReadWriter{
		Conn:    conn,
		scanner: bufio.NewScanner(conn),
	}
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
func (rw *ServerQueryReadWriter) ReadCommand() (rstr string, rerr error) {
	for {
		didScan := rw.scanner.Scan()
		if err := rw.scanner.Err(); err != nil {
			return "", err
		}
		if !didScan {
			return "", io.EOF
		}

		resultString := rw.scanner.Text()
		resultString = invalidCharRegex.ReplaceAllString(resultString, "")
		return resultString, nil
	}
}
