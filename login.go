package tclient

// Login is a simple wrapper for login/password auth.
// All this func do is waits for specific login/password prompts and sends login/password.
// You can always implement login manually.
func (c *TelnetClient) Login(login string, password string) (string, error) {
	// wait for login
	result := ""

	out, err := c.ReadUntil(c.loginPrompt)
	result = out
	if err != nil {
		return result, err
	}

	err = c.Write([]byte(login))
	if err != nil {
		return result, err
	}

	// and for password
	out, err = c.ReadUntil(c.passwordPrompt)
	result += out
	if err != nil {
		return result, err
	}

	err = c.Write([]byte(password))
	if err != nil {
		return result, err
	}

	// and wait for prompt
	out, err = c.ReadUntil(c.Prompt)
	result += out
	if err != nil {
		return result, err
	}

	return result, nil
}
