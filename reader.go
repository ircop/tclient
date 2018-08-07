package tclient

import (
	"time"
	"fmt"
	"bytes"
	"github.com/pkg/errors"
	"regexp"
)

func (c *TelnetClient) ReadUntil(waitfor string) (string, error) {
	var err error

	var b byte
	var prev byte
	var buf bytes.Buffer
	var lastLine bytes.Buffer
	if waitfor == "" {
		return buf.String(), fmt.Errorf(`Empty "waitfor" string given`)
	}

	rePrompt, err := regexp.Compile(waitfor)
	if err != nil {
		return buf.String(), fmt.Errorf(`Cannot compile "waitfor" regexp`)
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
			return buf.String() + lastLine.String(), fmt.Errorf("Operation timeout reached during read")
		default:
			prev = b
			b, err = c.readByte()
			if err != nil {
				return buf.String() + lastLine.String(), errors.Wrap(err, "Error during read")
			}

			// catch escape sequences
			if b == TELNET_IAC {
				seq := []byte{b}

				b2, err := c.readByte()
				if err != nil {
					return buf.String(), errors.Wrap(err, "Error while reading escape sequence")
				}

				seq = append(seq, b2)
				if b2 == TELNET_SB { // subnegotiation
					// read all until subneg. end.
					for {
						bn, err := c.readByte()
						if err != nil {
							return buf.String(), errors.Wrap(err, "Error while reading escape subnegotiation sequence")
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
						return buf.String(), errors.Wrap(err, "Error while reading IAC sequence")
					}
					seq = append(seq, bn)
				}

				// Sequence finished, do something with it:
				err = c.negotiate(seq)
				if err != nil {
					return buf.String(), errors.Wrap(err, "Failed to negotiate connection")
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
			// UPDATE: wirte to buf line by line, not every char. This is because we need to skip some lines (like pagination callbacks)
			buf.Write([]byte{b})
			//fmt.Printf("%v\t|\t%s\n",b,string(b))

			// check for regex matching
			if len(c.patterns) > 0 {
				for i := range c.patterns {
					if c.patterns[i].Re.Match(lastLine.Bytes()) {
						c.patterns[i].Cb()
						lastLine.Reset()
					}
				}
			}

			// check for CRLF.
			// We need last line to compare with prompt.
			if b == '\n' && prev == '\r' {
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
				return buf.String(), nil
			}
		}
	}
}
/*
func (c *TelnetClient) cutEscapes(buffer bytes.Buffer) string {
	bts := buffer.Bytes()
	result := ""

	inSequence := false
	for i := range bts {
		if bts[i] == 27 {
			inSequence = true
			continue
		}

		if inSequence {
			// 2) 0-?, @-~, ' ' - / === 48-63, 32-47, finish with 64-126
			if bts[i] == 91 {
				continue
			}
			if bts[i] >= 32 && bts[i] <= 63 {
				// just skip it
				continue
			}
			if bts[i] >= 64 && bts[i] <= 126 {
				// finish sequence
				inSequence = false
				continue
			}
		}

		result += string(bts[i])
	}

	return result
}*/

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
