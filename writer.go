package tclient

import (
	"time"
	"io"
	"github.com/pkg/errors"
)

// Write is the same as WriteRaw, but adds CRLF to given string
func (c *TelnetClient) Write(bytes []byte) error {
	bytes = append(bytes, '\r', '\n')

	return c.WriteRaw(bytes)
}

// WriteRaw writes raw bytes to tcp connection
func (c *TelnetClient) WriteRaw(bytes []byte) error {
	var wrote int
	var err error

	err = c.conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(c.Timeout)))
	if err != nil {
		return err
	}

	n, err := c.conn.Write(bytes)
	wrote += n

	if err != nil && err != io.ErrShortWrite {
		return errors.Wrap(err, "Failed to WriteRaw()")
	}

	return nil
}
