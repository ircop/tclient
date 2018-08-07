package tclient

// Cmd sends command and returns output.
// Actually, it just sends given command to server and waits for default prompt, and you can du this stuff manually.
func (c *TelnetClient) Cmd(cmd string) (string, error) {
	err := c.Write([]byte(cmd))
	if err != nil {
		return "", err
	}

	return c.ReadUntil(c.Prompt)
}

