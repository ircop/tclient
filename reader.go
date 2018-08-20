package tclient

import (
	"time"
	"fmt"
	"bytes"
	"github.com/pkg/errors"
	"regexp"
)

// ReadUntil reads tcp stream from server until 'waitfor' regex matches.
// Returns gathered output and error, if any.
// Any escape sequences are cutted out during reading for providing clean output for parsing/reading.
func (c *TelnetClient) ReadUntil(waitfor string) (string, error) {
	var err error

	var b byte
	var prev byte
	c.buf.Reset()
	//var buf bytes.Buffer
	var lastLine bytes.Buffer
	if waitfor == "" {
		return c.buf.String(), fmt.Errorf(`Empty "waitfor" string given`)
	}

	rePrompt, err := regexp.Compile(waitfor)
	if err != nil {
		return c.buf.String(), fmt.Errorf(`Cannot compile "waitfor" regexp`)
	}

	// run reading cycle
	inSequence := false
	//skipCurLine := false
	globalTout := time.After(time.Second * time.Duration(c.TimeoutGlobal))
	for {
		select {
		//case <- c.ctx.Done():
		//	return c.cutEscapes(buf), nil
		case <- globalTout:
			return c.buf.String() + lastLine.String(), fmt.Errorf("Operation timeout reached during read")
		default:
			prev = b
			b, err = c.readByte()
			if err != nil {
				return c.buf.String() + lastLine.String(), errors.Wrap(err, "Error during read")
			}

			// catch escape sequences
			if b == TELNET_IAC {
				seq := []byte{b}

				b2, err := c.readByte()
				if err != nil {
					return c.buf.String(), errors.Wrap(err, "Error while reading escape sequence")
				}

				seq = append(seq, b2)
				if b2 == TELNET_SB { // subnegotiation
					// read all until subneg. end.
					for {
						bn, err := c.readByte()
						if err != nil {
							return c.buf.String(), errors.Wrap(err, "Error while reading escape subnegotiation sequence")
						}
						seq = append(seq, bn)
						if bn == TELNET_SE {
							break
						}
					}
				} else {
					// not subsequence.
					bn, err := c.readByte()
					if err != nil {
						return c.buf.String(), errors.Wrap(err, "Error while reading IAC sequence")
					}
					seq = append(seq, bn)
				}

				// Sequence finished, do something with it:
				err = c.negotiate(seq)
				if err != nil {
					return c.buf.String(), errors.Wrap(err, "Failed to negotiate connection")
					//c.errChan <- fmt.Sprintf("Failed to negotiate connection: %s", err.Error())
				}
			}

			// cut out escape sequences
			if b == 27 {
				inSequence = true
				continue
			}
			if inSequence {
				// 2) 0-?, @-~, ' ' - / === 48-63, 32-47, finish with 64-126
				if b == 91 {
					continue
				}
				if b >= 32 && b <= 63 {
					// just skip it
					continue
				}
				if b >= 64 && b <= 126 {
					// finish sequence
					inSequence = false
					continue
				}
			}

			// this is not escape sequence, so write this byte to buffer
			// update: strip '\r'
			/*if b != '\r' {
				c.buf.Write([]byte{b})
			}*/

			// check for regex matching. Execute callback if matched.
			if len(c.patterns) > 0 {
				for i := range c.patterns {
					if c.patterns[i].Re.Match(lastLine.Bytes()) {
						c.patterns[i].Cb()
						lastLine.Reset()
						c.buf.Write([]byte{'\n'})
					}
				}
			}

			// remove \r ; remove backspaces
			if b == '\r' || b == 8 {
				continue
			}
			// check for CRLF.
			// We need last line to compare with prompt.
			if b == '\n' && prev == '\r' {
				lastLine.Write([]byte{b})
				c.buf.Write(lastLine.Bytes())
				lastLine.Reset()
			} else {
				lastLine.Write([]byte{b})
			}

			// After reading, we should check for regexp every time.
			// Unfortunately, we cant wait only CRLF, because prompt usually comes without CRLF.
			if rePrompt.Match(lastLine.Bytes()) {
				// we've catched required prompt.
				// now remove all escape sequences and return result
				//return buf.String() + lastLine.String(), nil
				return c.buf.String(), nil
			}
		}
	}
}


// read one byte from tcp stream
func (c *TelnetClient) readByte() (byte, error) {
	var err error
	var buffer [1]byte
	p := buffer[:]

	err = c.conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(c.Timeout)))
	if err != nil {
		return p[0], err
	}

	_, err = c.conn.Read(p)
	// error during read
	if err != nil {
		return p[0], errors.Wrap(err, "Error during readByte")
	}

	return p[0], nil
}
