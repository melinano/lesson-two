package dev03

import (
	"reflect"
	"testing"
)

func TestSortDataByColumn(t *testing.T) {
	// Sample input data
	lines := []string{
		"alice 25 October 50k apples",
		"bob 12 February 100 bananas",
		"carol 5 July 1m oranges",
		"david 18 April 500k bananas",
		"eve 30 December -10k grapes",
		"frank 8 August 2000 watermelons",
		"gary 12 February 100 bananas",
		"helen 5 July 1m oranges",
		"ivan 25 October 50k apples",
		"jack 18 April 500k bananas",
	}

	// Create a copy of the input data to preserve the original order for comparison
	originalLines := make([]string, len(lines))
	copy(originalLines, lines)

	// Expected result after sorting
	expectedResult := []string{
		"alice 25 October 50k apples",
		"ivan 25 October 50k apples",
		"bob 12 February 100 bananas",
		"david 18 April 500k bananas",
		"gary 12 February 100 bananas",
		"jack 18 April 500k bananas",
		"eve 30 December -10k grapes",
		"carol 5 July 1m oranges",
		"helen 5 July 1m oranges",
		"frank 8 August 2000 watermelons",
	}

	// Sort the data
	sortData(lines, 5, false, false, false, true, false)

	// Check if the sorted result matches the expected result
	if !reflect.DeepEqual(lines, expectedResult) {
		t.Errorf("SortData() failed, expected: %v, got: %v", expectedResult, lines)
	}
}

func TestSortDataByNumericWithSuffix(t *testing.T) {
	// Sample input data
	lines := []string{
		"vasya 14 january 10k aa",
		"jack 26 september 20000 ccc",
		"dimya 16 February 20000 ccc",
		"petya 13 MARCH 1m",
		"galya 18 April 200000k bbb",
		"denis 20 may -12k",
		"galya 8 October 200000k bbb",
		"chuck 2 February 20000 ccc",
	}

	// Expected result after sorting by numeric value with suffixes
	expectedResult := []string{
		"denis 20 may -12k",
		"vasya 14 january 10k aa",
		"jack 26 september 20000 ccc",
		"dimya 16 February 20000 ccc",
		"chuck 2 February 20000 ccc",
		"petya 13 MARCH 1m",
		"galya 18 April 200000k bbb",
		"galya 8 October 200000k bbb",
	}

	// Sort the data by numeric value with suffixes
	sortData(lines, 4, false, false, false, false, true)

	// Check if the sorted result matches the expected result
	if !reflect.DeepEqual(lines, expectedResult) {
		t.Errorf("SortData() failed for sorting by numeric value with suffixes, expected: %v, got: %v", expectedResult, lines)
	}
}

func TestSortDataByMonthName(t *testing.T) {
	// Sample input data
	lines := []string{
		"denis 20 may -12k",
		"vasya 14 january 10k aa",
		"chuck 16 March 3400 ccc",
		"gale 30 september 200000k bbb",
		"dimya 16 February 20000 ccc",
		"petya 13 MARCH 1m",
		"galya 18 April 200000k bbb",
		"dimya 16 February 20000 ccc",
	}

	// Expected result after sorting by month name
	expectedResult := []string{
		"vasya 14 january 10k aa",
		"dimya 16 February 20000 ccc",
		"dimya 16 February 20000 ccc",
		"chuck 16 March 3400 ccc",
		"petya 13 MARCH 1m",
		"galya 18 April 200000k bbb",
		"denis 20 may -12k",
		"gale 30 september 200000k bbb",
	}

	// Sort the data by month name
	sortData(lines, 3, false, false, true, false, false)

	// Check if the sorted result matches the expected result
	if !reflect.DeepEqual(lines, expectedResult) {
		t.Errorf("SortData() failed for sorting by month name, expected: %v, got: %v", expectedResult, lines)
	}
}
