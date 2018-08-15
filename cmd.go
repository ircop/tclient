package tclient

// Cmd sends command and returns output
func (c *TelnetClient) Cmd(cmd string) (string, error) {
	err := c.Write([]byte(cmd))
	if err != nil {
		return "", err
	}

	return c.ReadUntil(c.prompt)
}

