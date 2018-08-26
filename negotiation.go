package tclient

// ignore: TELNET_AYT, TELNET_AO
// accept: 0x01, 0x03, 0x18, 0x1f == TELOPT_ECHO, TELOPT_SGA, TELOPT_TTYPE, TELOPT_NAWS
func (c *TelnetClient) negotiate(sequence []byte) error {
	if len(sequence) < 3 {
		// do nothing
		return nil
	}
	//fmt.Println(sequence)
	var err error

	// DO sequence
	if sequence[1] == TELNET_DO && len(sequence) == 3 {
		switch sequence[2] {
		case TELOPT_TTYPE:
			err = c.WriteRaw([]byte{TELNET_IAC, TELNET_WILL, TELOPT_TTYPE})
			break
		case TELOPT_SB_NEV_ENVIRON:
			// iac do newinv - iac will newinv
			err = c.WriteRaw([]byte{TELNET_IAC, TELNET_WILL, TELOPT_SB_NEV_ENVIRON})
			break
		case TELOPT_NAWS:
			// wont naws
			err = c.WriteRaw([]byte{TELNET_IAC, TELNET_WONT, TELOPT_NAWS})
			break
		default:
			// accept any other 'do'
			err = c.WriteRaw([]byte{TELNET_IAC, TELNET_WILL, sequence[2]})
			break
		}
	}

	// WONT -> DONT
	if sequence[1] == TELNET_WONT && len(sequence) == 3 {
		err = c.WriteRaw([]byte{TELNET_IAC, TELNET_DONT, sequence[2]})
	}

	// WILL -> DO
	if sequence[1] == TELNET_WILL && len(sequence) == 3 {
		err = c.WriteRaw([]byte{TELNET_IAC, TELNET_DO, sequence[2]})
	}

	// subseq SEND request
	if len(sequence) == 6 && sequence[1] == TELNET_SB && sequence[3] == TELOPT_SB_SEND {
		// what to send?
		switch(sequence[2]) {
		case TELOPT_TTYPE:
			// set terminal to xterm
			err = c.WriteRaw([]byte{TELNET_IAC, TELNET_SB, TELOPT_TTYPE, TELOPT_SB_IS, 'X', 'T', 'E', 'R', 'M', TELNET_IAC, TELNET_SE})
			break
		case TELOPT_SB_NEV_ENVIRON:
			// send new-env -> is new env
			err = c.WriteRaw([]byte{TELNET_IAC, TELNET_SB, TELOPT_SB_NEV_ENVIRON, TELOPT_SB_IS, TELNET_IAC, TELNET_SE})
			break
		default:
			// accept all
			err = c.WriteRaw([]byte{TELNET_IAC, TELNET_SB, sequence[2], 0, TELNET_IAC, TELNET_SE})
			break
		}
	}

	// subseq IS request
	if len(sequence) == 6 && sequence[1] == TELNET_SB && sequence[3] == TELOPT_SB_IS {
		//
	}

	if err != nil {
		return err
	}
	return nil
}
