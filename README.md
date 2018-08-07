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

	// you can omit this, or do auth stuff manually by calling `readUntil` with login/password prompts
	out, err := client.Login("script2", "wre4fel")
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




# TODO

Implement pagination parsing/manipulating (various network devices paging their output)
