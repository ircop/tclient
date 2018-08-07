# tclient

Simple telnet client lib, written in golang.

Example usage:

main.go:
```
	client := tclient.New(5, "")
	err := client.Open("10.10.10.10", 23)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// you can omit this, or do auth stuff manually by calling `ReadUntil` with login/password prompts
	out, err := client.Login("script2", "pw3")
	if err != nil {
		panic(err)
	}
	fmt.Printf(out)

	out, err = client.Cmd("show time")
	if err != nil {
		panic(err)
	}
	fmt.Printf(out)
```

Output: 

![Output](https://i.imgur.com/2M91MEN.png)


# Matching callbacks

You can define regular expressions and callbacks that would be called when current output string will match one of regexps.

For example, we need to catch pagination on D-Link switch. Sample `show switch` paginated output:

![Show switch](https://i.imgur.com/PoUBDyQ.png)

There is no prompt, so app will stuck on this. So we need to catch something like `CTRL+C ESC q Quit SPACE n Next Page ENTER Next Entry a All` and send 'n' (next page) or 'a' (all pages). Like this:

```
	// matching "CTRL+C ESC q Quit SPACE n Next Page ENTER Next Entry a All"
	err = client.RegisterCallback(`(?msi:CTRL\+C.+?a A[Ll][Ll]\s*)`, func() {
		client.WriteRaw([]byte("a"))
	})
	if err != nil {
		panic(err)
	}
```

Note we are using WriteRaw(), not Write(), because Write() adds CRLF to given string.



# TODO

~~Implement pagination parsing/manipulating (various network devices paging their output)~~

~~Implement parsing regexps with callbacks. For output pagination purposes, etc.~~
