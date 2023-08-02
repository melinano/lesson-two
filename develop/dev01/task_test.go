package dev01

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

// Testing the proper output of the error code and message to stderr
func TestPrintCurrentTimeError(t *testing.T) {
	fmt.Println("Testing for proper error output")
	// Create a mock NTP server
	ntpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate an error by responding with an empty body
	}))
	defer ntpServer.Close()

	// Replace the NTP server URL with the mock server URL
	ntpServerURL = ntpServer.URL

	if os.Getenv("FLAG") == "1" {
		PrintCurrentTime()
		return
	}

	cmd := exec.Command(os.Args[0], "test.run=TestPrintCurrentTime")
	cmd.Env = append(os.Environ(), "FLAG=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with error %v, want exit status 1", err)
}
