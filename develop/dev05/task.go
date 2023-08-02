package dev05

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// Flags struct to hold the values of different command line flags
type Flags struct {
	afterLines   int
	beforeLines  int
	contextLines int
	count        bool
	ignoreCase   bool
	invert       bool
	fixed        bool
	lineNumber   bool
}

// initFlags initializes the default values of the command line flags
func initFlags(flags *Flags) {
	flag.IntVar(&flags.afterLines, "A", 0, "Print +N lines after a match")
	flag.IntVar(&flags.beforeLines, "B", 0, "Print +N lines before a match")
	flag.IntVar(&flags.contextLines, "C", 0, "Print ±N lines around a match")
	flag.BoolVar(&flags.count, "c", false, "Print number of lines")
	flag.BoolVar(&flags.ignoreCase, "i", false, "Ignore case when matching")
	flag.BoolVar(&flags.invert, "v", false, "Invert the match")
	flag.BoolVar(&flags.fixed, "F", false, "Match exact line")
	flag.BoolVar(&flags.lineNumber, "n", false, "Print line numbers")
}

// main function to run the filtering utility
func main() {
	var flags Flags
	initFlags(&flags)
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: grep [flags] pattern [file]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	pattern := args[0]
	filePath := ""
	if len(args) > 1 {
		filePath = args[1]
	}

	// Compile the regular expression based on the pattern and flags
	regex := pattern
	if flags.fixed {
		regex = regexp.QuoteMeta(pattern)
	}
	if flags.ignoreCase {
		regex = "(?i)" + regex
	}
	re, err := regexp.Compile(regex)
	if err != nil {
		fmt.Println("Error compiling regular expression:", err)
		os.Exit(1)
	}

	var reader io.Reader
	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	// Start reading and filtering the lines from the input
	scanner := bufio.NewScanner(reader)
	var (
		matchingLines int
		outputLines   []string
		lineNumber    int
		afterLines    int
		beforeLines   int
	)

	if flags.contextLines > 0 {
		afterLines = flags.contextLines
		beforeLines = flags.contextLines
	} else {
		afterLines = flags.afterLines
		beforeLines = flags.beforeLines
	}

	// Initialize a circular queue to store lines before the match
	beforeBuffer := make([]string, beforeLines)
	// Index to keep track of the position in the circular queue
	index := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNumber++

		// Check for a match
		matched := re.MatchString(line)
		if (matched && !flags.invert) || (!matched && flags.invert) {
			matchingLines++

			// Print lines before the match
			for i := 0; i < beforeLines; i++ {
				if beforeBuffer[index] != "" {
					matchingLines++
					outputLines = append(outputLines, formatOutputLine(flags, lineNumber-beforeLines+i, beforeBuffer[index]))
				}
				index = (index + 1) % beforeLines
			}

			// Print the matching line
			outputLines = append(outputLines, formatOutputLine(flags, lineNumber, line))

			afterLinesCounter := afterLines
			// Print lines after the match
			for afterLinesCounter > 0 && scanner.Scan() {
				lineNumber++
				afterLinesCounter--
				line = scanner.Text()
				outputLines = append(outputLines, formatOutputLine(flags, lineNumber, line))
			}

			// Clear the beforeBuffer after the match
			for i := 0; i < beforeLines; i++ {
				beforeBuffer[index] = ""
				index = (index + 1) % beforeLines
			}
		} else {
			if len(beforeBuffer) < 3 {
				beforeBuffer = append(beforeBuffer, line)
			} else {
				// Add the line to the beforeBuffer
				beforeBuffer[index] = line
				index = (index + 1) % beforeLines
			}
		}
	}

	// Print the final output
	if flags.count {
		fmt.Println(matchingLines)
	} else {
		for _, line := range outputLines {
			fmt.Println(line)
		}
	}
}

// formatOutputLine formats the output line with line numbers if required
func formatOutputLine(flags Flags, lineNumber int, line string) string {
	if flags.lineNumber {
		return fmt.Sprintf("%d:%s", lineNumber, line)
	}
	return line
}
