package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

/*
Реализовать простейший telnet-клиент.

Примеры вызовов:
	go-telnet --timeout=10s host port go-telnet mysite.ru 8080
	go-telnet --timeout=3s 1.1.1.1 123


Требования:
	- 	Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP.
		После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
	- 	Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)
	- 	При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера,
		программа должна также завершаться. При подключении к несуществующему сервер, программа должна завершаться через timeout
*/

type TelnetClient struct {
	remoteAddr string
	timeout    time.Duration
	in         io.Reader
	out        io.Writer
}

// NewTelnetClient creates a new Telnet client
func NewTelnetClient(addr string, timeout time.Duration, in io.Reader, out io.Writer) *TelnetClient {
	return &TelnetClient{
		remoteAddr: addr,
		timeout:    timeout,
		in:         in,
		out:        out,
	}
}

// Connect connects to the remote server and starts copying data between the server and the client
func (t *TelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.remoteAddr, t.timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to remote addr: %v", err)
	}
	defer conn.Close()

	// Start a goroutine to read from the connection and write to the output
	go func() {
		if _, err := io.Copy(t.out, conn); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	// Copy data from the input to the connection
	if _, err := io.Copy(conn, t.in); err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}

	return nil
}

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "timeout")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go-telnet --timeout=10s host port")
		os.Exit(1)
	}

	addr := net.JoinHostPort(args[0], args[1])
	telnetClient := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)

	if err := telnetClient.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
