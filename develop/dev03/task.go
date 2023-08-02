package dev03

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
=== Утилита sort ===

Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры):
на входе подается файл из несортированных строк, на выходе — файл с отсортированными.

Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// Command-line flags for sorting options
	var sortByColumn int
	var sortByNumeric bool
	var reverseSort bool
	var uniqueRows bool
	var sortByMonthName bool
	var ignoreTrailingSpaces bool
	var checkSorted bool
	var sortByNumericWithSuffix bool

	flag.IntVar(&sortByColumn, "k", 0, "Specify column for sorting (1-based index)")
	flag.BoolVar(&sortByNumeric, "n", false, "Sort by numeric value")
	flag.BoolVar(&reverseSort, "r", false, "Sort in reverse order")
	flag.BoolVar(&uniqueRows, "u", false, "Do not output repeated rows")
	flag.BoolVar(&sortByMonthName, "M", false, "Sort by month name")
	flag.BoolVar(&ignoreTrailingSpaces, "b", false, "Ignore tail spaces")
	flag.BoolVar(&checkSorted, "c", false, "Check if the data is sorted")
	flag.BoolVar(&sortByNumericWithSuffix, "h", false, "Sort by numeric value including suffixes")

	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Println("Usage: go run main.go [flags] input_file output_file")
		return
	}

	inputFile := flag.Arg(0)
	outputFile := flag.Arg(1)

	lines, err := readLines(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	if checkSorted {
		if isSorted(lines, sortByColumn, sortByNumeric, sortByMonthName, sortByNumericWithSuffix) {
			fmt.Println("Data is sorted.")
		} else {
			fmt.Println("Data is not sorted.")
		}
		return
	}

	sortData(lines, sortByColumn, sortByNumeric, reverseSort, sortByMonthName, ignoreTrailingSpaces, sortByNumericWithSuffix)

	if uniqueRows {
		lines = removeDuplicates(lines)
	}

	if err := writeLines(outputFile, lines); err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	fmt.Println("Sorting successful. Result written to", outputFile)
}

// readLines reads lines from the specified file and returns them as a slice of strings.
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// writeLines writes the lines to the specified file.
func writeLines(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	return writer.Flush()
}

// sortData sorts the lines based on the specified sorting options.
func sortData(lines []string, sortByColumn int, sortByNumeric, reverseSort, sortByMonthName, ignoreTrailingSpaces, sortByNumericWithSuffix bool) {
	// Define the sorting logic using SliceStable
	sort.SliceStable(lines, func(i, j int) bool {
		lineA := lines[i]
		lineB := lines[j]

		// variable that says if i comes before j
		var iComesFirst bool

		// Default behavior: Compare the strings in the specified column
		fieldA, fieldB := getField(lineA, sortByColumn), getField(lineB, sortByColumn)

		// Handle ignoring trailing spaces
		if ignoreTrailingSpaces {
			fieldA = strings.TrimSpace(fieldA)
			fieldB = strings.TrimSpace(fieldB)
		}

		// Sort entire lines if sortByColumn is 0
		if sortByColumn == 0 {
			fieldA = lineA
			fieldB = lineB
		}

		switch {
		// Handle sorting by month name
		case sortByMonthName:
			monthA, okA := parseMonthName(lineA, sortByColumn)
			monthB, okB := parseMonthName(lineB, sortByColumn)

			if okA && okB {
				iComesFirst = monthA < monthB
			} else if okA {
				iComesFirst = true
			} else if okB {
				iComesFirst = false
			}

		// Handle sorting by numeric value with suffixes
		case sortByNumericWithSuffix:
			numA, okA := parseNumericWithSuffix(lineA, sortByColumn)
			numB, okB := parseNumericWithSuffix(lineB, sortByColumn)

			if okA && okB {
				iComesFirst = numA < numB
			} else if okA {
				iComesFirst = true
			} else if okB {
				iComesFirst = false
			}
		// Handle sorting by numeric value
		case sortByNumeric:
			numA, okA := parseFloat(fieldA)
			numB, okB := parseFloat(fieldB)

			if okA && okB {
				iComesFirst = numA < numB
			} else if okA {
				iComesFirst = true
			} else if okB {
				iComesFirst = false
			}

		default:
			iComesFirst = fieldA < fieldB
		}

		// Handle reverse sorting
		if reverseSort {
			return !iComesFirst
		}

		return iComesFirst

	})
}

// getField retrieves the field from the line based on the specified column index.
func getField(line string, column int) string {
	if column <= 0 {
		return ""
	}

	fields := strings.Fields(line)
	if column <= len(fields) {
		return fields[column-1]
	}
	return ""
}

// parseFloat attempts to parse the given string as a float64.
func parseFloat(s string) (float64, bool) {
	num, err := strconv.ParseFloat(s, 64)
	return num, err == nil
}

// parseMonthName attempts to parse the month name from the specified column.
func parseMonthName(line string, column int) (time.Month, bool) {
	monthStr := getField(line, column)
	for i := 1; i <= 12; i++ {
		if strings.EqualFold(monthStr, time.Month(i).String()) {
			return time.Month(i), true
		}
	}
	return 0, false
}

// parseNumericWithSuffix attempts to parse the numeric value with suffix from the specified column.
func parseNumericWithSuffix(line string, column int) (float64, bool) {
	field := getField(line, column)
	if len(field) == 0 {
		return 0, false
	}

	lastChar := field[len(field)-1]
	if lastChar >= '0' && lastChar <= '9' {
		return parseFloat(field)
	}

	suffix := field[len(field)-1]
	numPart := field[:len(field)-1]
	multiplier := 1.0

	switch suffix {
	case 'k', 'K':
		multiplier = 1000.0
	case 'm', 'M':
		multiplier = 1000000.0
	case 'b', 'B':
		multiplier = 1000000000.0
	case 't', 'T':
		multiplier = 1000000000000.0
	default:
		return 0, false
	}

	num, _ := parseFloat(numPart)

	return num * multiplier, true
}

// removeDuplicates removes duplicate lines from the input slice.
func removeDuplicates(lines []string) []string {
	uniqueMap := make(map[string]struct{})
	var uniqueLines []string

	for _, line := range lines {
		// Ignore leading and trailing whitespaces for duplicate comparison
		trimmedLine := strings.TrimSpace(line)

		if _, ok := uniqueMap[trimmedLine]; !ok {
			uniqueMap[trimmedLine] = struct{}{}
			uniqueLines = append(uniqueLines, line)
		}
	}

	return uniqueLines
}

// isSorted checks if the lines are sorted based on the specified sorting options.
func isSorted(lines []string, sortByColumn int, sortByNumeric, sortByMonthName, sortByNumericWithSuffix bool) bool {
	for i := 1; i < len(lines); i++ {
		lineA := lines[i-1]
		lineB := lines[i]

		if sortByMonthName {
			monthA, okA := parseMonthName(lineA, sortByColumn)
			monthB, okB := parseMonthName(lineB, sortByColumn)

			if okA && okB && monthA > monthB {
				return false
			} else if okA && !okB {
				return false
			} else if !okA && okB {
				return true
			}
		}

		if sortByNumericWithSuffix {
			numA, okA := parseNumericWithSuffix(lineA, sortByColumn)
			numB, okB := parseNumericWithSuffix(lineB, sortByColumn)

			if okA && okB && numA > numB {
				return false
			} else if okA && !okB {
				return false
			} else if !okA && okB {
				return true
			}
		}

		fieldA, fieldB := getField(lineA, sortByColumn), getField(lineB, sortByColumn)

		if sortByNumeric {
			numA, okA := parseFloat(fieldA)
			numB, okB := parseFloat(fieldB)

			if okA && okB && numA > numB {
				return false
			} else if okA && !okB {
				return false
			} else if !okA && okB {
				return true
			}
		}

		if fieldA > fieldB {
			return false
		}
	}

	return true
}
