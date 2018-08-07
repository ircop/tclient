package tclient

import (
	"testing"
	"net"
	"github.com/pkg/errors"
	"time"
)

func TestNew(t *testing.T) {
	c := New(-1, "")
	if c.Timeout != 1 {
		t.Fatal("Negative timeout should be modified to 1")
	}
}

func TestDialFailed(t *testing.T) {
	c := New(0,"")
	err := c.Open("127.0.0.1", 12345)
	if err == nil {
		t.Fatal("127.0.0.1:12345 should fails")
	}
}

func TestDialOk(t *testing.T) {
	l, err := net.Listen("tcp4", ":33333")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	c := New(0, "")
	err = c.Open("127.0.0.1", 33333)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Failed to dial localhost"))
	}
}

func TestReadUntilSimple(t *testing.T) {
	s, c := net.Pipe()
	defer s.Close()
	defer c.Close()
	client := New(0, "")
	client.conn = c

	contents := `asd qwe zxc 123 >`
	go func() {
		_, e := s.Write([]byte(contents))
		if e != nil {
			t.Fatal(e, "Failed to write to pipe")
		}
	}()

	out, err := client.ReadUntil(`>$`)
	if err != nil {
		t.Fatal(errors.Wrap(err, "Error during pipe read"))
	}

	if out != contents {
		t.Fatal("Recorded and read data doesn't match!")
	}
}

func TestREadUntilTimeout(t *testing.T) {
	s, c := net.Pipe()
	defer s.Close()
	defer c.Close()
	client := New(0, "")
	client.conn = c

	_, err := client.ReadUntil(`>$`)
	if err == nil {
		t.Fatal("Should timeout")
	}
}

func TestComplexLogin(t *testing.T) {
	login := "testLogin"
	pw := "testPw"

	s, c := net.Pipe()
	defer s.Close()
	defer c.Close()
	client := New(0, "")
	client.conn = c

	preLogin := `
              DES-3028 Fast Ethernet Switch Command Line Interface

                            Firmware: Build 2.70.B06
           Copyright(C) 2008 D-Link Corporation. All rights reserved.
UserName:`
	prePw := `PassWord:`
	prePrompt := `
DES-3028:5#`

	// server should also read messages from client...
	go func() {
		for {
			s.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(999)))
			var buffer [1]byte
			p := buffer[:]
			_, err := s.Read(p)
			if err != nil {
				return
			}
		}
	}()

	go func () {
		_, e := s.Write([]byte(preLogin))
		if e != nil {
			t.Fatalf("Error writing preLogin: %s", e.Error())
		}
	}()

	_, err := client.ReadUntil(client.loginPrompt)
	if err != nil {
		t.Fatalf("Error reading login prompt: %s", err.Error())
	}

	err = client.Write([]byte(login))
	if err != nil {
		t.Fatalf("error writing login: %s", err.Error())
	}

	go func () {
		_, e := s.Write([]byte(prePw))
		if e != nil {
			t.Fatalf("Error writing preLogin: %s", e.Error())
		}
	}()

	err = client.Write([]byte(pw))
	if err != nil {
		t.Fatalf("error writing password: %s", err.Error())
	}

	_, err = client.ReadUntil(client.passwordPrompt)
	if err != nil {
		t.Fatalf("Error reading password prompt: %s", err.Error())
	}

	go func () {
		_, e := s.Write([]byte(prePrompt))
		if e != nil {
			t.Fatalf("Error writing prompt: %s", e.Error())
		}
	}()

	out, err := client.ReadUntil(client.Prompt)
	if err != nil {
		t.Fatalf("Error reading prompt: %s\nlast: %s", err.Error(), out)
	}
}

