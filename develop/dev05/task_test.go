package dev05

import (
	"os/exec"
	"strings"
	"testing"
)

func TestGrep(t *testing.T) {

	filePath := "shakespeare.txt"

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "TestCase1",
			args:     []string{"-n", "-A", "1", "beauty", filePath},
			expected: "9:That thereby beauty’s rose might never die,\n10:But as the riper should by time decease,\n11:His tender heir might bear his memory:\n",
		},
		{
			name:     "TestCase2",
			args:     []string{"-n", "-B", "2", "small", filePath},
			expected: "13:Thy youth’s proud livery so gazed on now,\n14:Will be a tattered weed of small worth held:\n",
		},
		{
			name:     "TestCase3",
			args:     []string{"-n", "-C", "2", "world", filePath},
			expected: "9:That thereby beauty’s rose might never die,\n10:But as the riper should by time decease,\n11:His tender heir might bear his memory:\n12:But thou contracted to thine own bright eyes,\n13:Feed’st thy light’s flame with self-substantial fuel,\n",
		},
		{
			name:     "TestCase4",
			args:     []string{"-n", "-i", "mAkIng", filePath},
			expected: "9:Making a famine where abundance lies,\n",
		},
		{
			name:     "TestCase5",
			args:     []string{"-n", "-v", "glutton", filePath},
			expected: "12:  Pity the world, or else this glutton be,\n",
		},
	}

	// run the tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd := exec.Command("./task.go", test.args...)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Command execution failed: %v", err)
			}

			actual := string(output)
			if !strings.Contains(actual, test.expected) {
				t.Errorf("Expected:\n%s\n\nActual:\n%s", test.expected, actual)
			}
		})
	}
}
