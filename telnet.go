// Package tclient provides simple telnet client library, mainly designed for iteractions with
// various network equipment such as switches, routers, etc.
package tclient

import (
	"fmt"
	"net"
	"time"
	"regexp"
	"bytes"
)

type callbackPattern struct {
	Re			*regexp.Regexp
	Cb			func()
	SkipLine	bool
}

// TelnetClient is telnet client struct itself
type TelnetClient struct {
	// Timeout of read/write operations
	Timeout			int
	// TimeoutGlobal is global operation timeout; i.e. like stucked in DLink-like refreshing pagination
	TimeoutGlobal 	int
	login			string
	password		string
	prompt        	string
	conn          	net.Conn
	closed        	bool
	Options       	[]int

	loginPrompt		string
	passwordPrompt	string

	buf				bytes.Buffer

	patterns		[]callbackPattern
}

// New func creates new TelnetClient instance.
// First argument is network (r/w) timeout, second is prompt. Default prompt is "(?msi:[\$%#>]$)"
func New(tout int, login string, password string, prompt string) *TelnetClient {
	if tout < 1 {
		tout = 1
	}
	c := TelnetClient{
		Timeout:        tout,
		login:			login,
		password:		password,
		prompt:         `(?msi:[\$%#>]$)`,
		loginPrompt:    `[Uu]ser[Nn]ame\:$`,
		passwordPrompt: `[Pp]ass[Ww]ord\:$`,
		closed:         false,
		Options:        make([]int,0),
	}

	if prompt != "" {
		c.prompt = prompt
	}

	// Global timeout defaults to 3 * rw timeout
	c.TimeoutGlobal = c.Timeout * 2

	// set default options
	// we will accept an offer from remote sidie for it to echo and suppress goaheads
	c.SetOpts([]int{TELOPT_ECHO, TELOPT_SGA})

	return &c
}

// GlobalTimeout sets timeout for app operations, where net.Conn deadline could not be useful.
// For example stucking in pagination, while some network devices refreshing their telnet screen - so
// we cannot reach read timeout.
func (c *TelnetClient) GlobalTimeout(t int) {
	c.TimeoutGlobal = t
}

// You may need to change password for enable (because prompt is same as login)
func (c *TelnetClient) SetPassword(pw string) {
	c.password = pw
}

// SetPrompt allows you to change prompt without re-creating ssh client
func (c *TelnetClient) SetPrompt(prompt string) {
	c.prompt = prompt
}

// SetLoginPrompt sets custom login prompt. Default is "[Uu]ser[Nn]ame\:$"
func (c *TelnetClient) SetLoginPrompt(s string) {
	c.loginPrompt = s
}

// SetPasswordPrompt sets custom password prompt. Default is "[Pp]ass[Ww]ord\:$"
func (c *TelnetClient) SetPasswordPrompt(s string) {
	c.passwordPrompt = s
}

// Close closes telnet connection. You can use it with defer.
func (c *TelnetClient) Close() {
	if !c.closed {
		c.conn.Close()
	}
	c.closed = false
}

// FlushOpts flushes default options set
func (c *TelnetClient) FlushOpts() {
	c.Options = make([]int, 0)
}

// SetOpts prepares default options set
func (c *TelnetClient) SetOpts(opts []int) error {
	for _, opt := range opts {
		if opt > 255 {
			return fmt.Errorf("Bad telnet option %d: option > 255", opt)
		}

		c.Options = append(c.Options, opt)
	}

	return nil
}

// Open tcp connection
func (c *TelnetClient) Open(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)

	var err error
	c.conn, err = net.DialTimeout("tcp", addr, time.Second * time.Duration(c.Timeout))
	if err != nil {
		c.closed = true
		return err
	}

	c.closed = false

	// and login right now
	_, err = c.Login(c.login, c.password)

	return err
}

// RegisterCallback registers new callback based on regex string. When current output string matches given
// regex, callback is called. Returns error if regex cannot be compiled.
func (c *TelnetClient) RegisterCallback(pattern string, callback func()) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	c.patterns = append(c.patterns, callbackPattern{
		Cb:callback,
		Re:re,
		})

	return nil
}

// GetBuffer returns current buffer from reader as a string
func (c *TelnetClient) GetBuffer() string {
	return c.buf.String()
}