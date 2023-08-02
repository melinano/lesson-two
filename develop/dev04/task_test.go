package dev04

import (
	"reflect"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	// Input data with anagrams in Russian
	tests := []struct {
		name  string
		words []string
		want  map[string][]string
	}{
		{
			name:  "nil words",
			words: nil,
			want:  make(map[string][]string),
		},
		{
			name:  "empty words",
			words: []string{},
			want:  make(map[string][]string),
		},
		{
			name:  "one word",
			words: []string{"дом"},
			want:  make(map[string][]string),
		},
		{
			name:  "base example",
			words: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			want: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "base example and one word",
			words: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "дом"},
			want: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "base example in different registers",
			words: []string{"пЯтАк", "ПяТкА", "тяпка", "ЛИСток", "слиток", "стоЛИК"},
			want: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:  "base example with double words",
			words: []string{"пятак", "пятак", "пятка", "тяпка", "листок", "слиток", "слиток", "слиток", "столик"},
			want: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"пятак":  {"пятак", "пятка", "тяпка"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindAnagrams(tt.words)

			// Check if the actual result matches the expected result
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAnagramSets() failed, expected: %v, got: %v", tt.want, got)
			}
		})
	}
}
