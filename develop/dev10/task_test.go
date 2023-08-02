package main

import (
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestTelnetClient_Connect(t *testing.T) {
	// Start a test server
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer ln.Close()

	go func(t *testing.T) {
		conn, err := ln.Accept()
		if err != nil {
			t.Errorf("failed to accept connection: %v", err)
		}
		defer conn.Close()

		fmt.Fprintln(conn, "Hello, client!")
	}(t)

	// Initialize a TelnetClient with a fake reader and a string builder to capture the output
	in := strings.NewReader("Hello, server!\n")
	out := &strings.Builder{}
	tc := NewTelnetClient(ln.Addr().String(), time.Second, in, out)

	// Connect to the test server and check the output
	if err := tc.Connect(); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	got := out.String()
	want := "Hello, client!\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
