package dev02

import "testing"

func TestUnpackString(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		err      error
	}{
		{"a4bc2d5e", "aaaabccddddde", nil},
		{"abcd", "abcd", nil},
		{"45", "44444", nil},
		{"", "", nil},
		{"qwe\\4\\5", "qwe45", nil},
		{"qwe\\45", "qwe44444", nil},
		{`qwe\\5`, `qwe\\\\\`, nil},
		{`qwe\\\5`, `qwe\5`, nil},
		{`qwe\\2\`, "", errInvalidEsc}, // Invalid escape, should not unpack the last '\'
		{"qwe\\", "", errInvalidEsc},   // Invalid escape at the end
		{`qwe\\1\\`, `qwe\\`, nil},     // Valid escape sequence with a single repetition
		{`qwe\\01\\`, `qwe1\`, nil},    // Valid escape sequence with a single repetition (leading zeros)
	}

	for _, testCase := range testCases {
		result, err := UnpackString(testCase.input)

		if testCase.err != nil && err == nil {
			t.Errorf("Expected error, but got nil for input: %q", testCase.input)
		}

		if testCase.err == nil && err != nil {
			t.Errorf("Unexpected error for input %q: %v", testCase.input, err)
		}

		if testCase.err != nil && err != nil && testCase.err.Error() != err.Error() {
			t.Errorf("Unexpected error for input %q. Expected: %v, Got: %v", testCase.input, testCase.err, err)
		}

		if result != testCase.expected {
			t.Errorf("Unexpected result for input %q. Expected: %q, Got: %q", testCase.input, testCase.expected, result)
		}
	}
}
