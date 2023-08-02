package main

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

func TestDownload(t *testing.T) {
	// Start the downloading process for the URL
	startURL := "http://example.com"
	queue := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go download(startURL, queue, &wg)

	go func() {
		wg.Wait()
		close(queue)
	}()

	for uri := range queue {
		if _, err := os.Stat(urlToFileName(uri)); os.IsNotExist(err) {
			wg.Add(1)
			go download(uri, queue, &wg)
		}
	}

	// Wait for the download to finish
	wg.Wait()

	// Compare the contents of the downloaded file and the expected file
	expectedContent, err := ioutil.ReadFile("example.html")
	if err != nil {
		t.Errorf("Failed to read example.html: %v", err)
	}

	downloadedContent, err := ioutil.ReadFile(urlToFileName(startURL))
	if err != nil {
		t.Errorf("Failed to read downloaded file: %v", err)
	}

	if string(downloadedContent) != string(expectedContent) {
		t.Errorf("Contents of downloaded file and example.html do not match")
	}
}
