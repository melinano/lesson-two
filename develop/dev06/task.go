package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var fields = flag.String("f", "", "List of fields to display")
var delimiter = flag.String("d", "\t", "Delimiter to use")
var separated = flag.Bool("s", false, "Only lines with the delimiter")

func main() {
	// Parse command line arguments
	flag.Parse()

	// Convert fields from comma-separated values to a slice
	f := strings.Split(*fields, ",")

	processInput(os.Stdin, os.Stdout, f, *delimiter, *separated)
}

func processInput(input io.Reader, output io.Writer, fields []string, delimiter string, separated bool) {
	scanner := bufio.NewScanner(input)

	// Iterate over lines from input
	for scanner.Scan() {
		line := scanner.Text()

		// If -s is set and line does not contain the delimiter, skip it
		if separated && !strings.Contains(line, delimiter) {
			continue
		}

		// Split line into fields by the delimiter
		columns := strings.Split(line, delimiter)

		// Iterate over requested fields and print them
		for _, field := range fields {
			index := convertToInt(field)
			if index >= 0 && index < len(columns) {
				fmt.Fprint(output, columns[index])
			}
			fmt.Fprint(output, delimiter)
		}
		fmt.Fprintln(output)
	}
}

func convertToInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return result
}
