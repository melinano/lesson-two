package main

import (
	"strings"
	"testing"
)

func TestConvertToInt(t *testing.T) {
	if convertToInt("1") != 1 {
		t.Errorf("convertToInt(\"1\") = %d; want 1", convertToInt("1"))
	}
	if convertToInt("abc") != -1 {
		t.Errorf("convertToInt(\"abc\") = %d; want -1", convertToInt("abc"))
	}
}

func TestProcessInput(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		output    string
		fields    []string
		delimiter string
		separated bool
	}{
		{
			"Basic case",
			"a\tb\tc\n",
			"a\t\n",
			[]string{"0"},
			"\t",
			false,
		},
		{
			"Multiple fields",
			"a\tb\tc\n",
			"a\tc\t\n",
			[]string{"0", "2"},
			"\t",
			false,
		},
		{
			"Nonexistent field",
			"a\tb\tc\n",
			"\t\n",
			[]string{"3"},
			"\t",
			false,
		},
		{
			"Separated option",
			"a\tb\tc\na\n",
			"a\t\n",
			[]string{"0"},
			"\t",
			true,
		},
		{
			"Different delimiter",
			"a,b,c\n",
			"a,\n",
			[]string{"0"},
			",",
			false,
		},
		{
			"Invalid field",
			"a\tb\tc\n",
			"\t\n",
			[]string{"a"},
			"\t",
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := strings.NewReader(tc.input)
			output := &strings.Builder{}

			processInput(input, output, tc.fields, tc.delimiter, tc.separated)

			if output.String() != tc.output {
				t.Errorf("got = %q; want %q", output.String(), tc.output)
			}
		})
	}
}
