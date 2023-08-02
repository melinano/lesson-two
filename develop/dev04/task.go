package dev04

import (
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func FindAnagrams(dictionary []string) map[string][]string {
	tempAnagramMap := make(map[string][]string)
	anagramMap := make(map[string][]string)

	for _, word := range dictionary {
		// convert to lowercase
		word = strings.ToLower(word)
		// a word, where the runes are sorted alphabetically serves as key for a temporary map
		sortedWord := sortString(word)
		if val, ok := tempAnagramMap[sortedWord]; ok {
			if contains(val, word) {
				continue
			}
			tempAnagramMap[sortedWord] = append(tempAnagramMap[sortedWord], word)
		} else {
			tempAnagramMap[sortedWord] = []string{word}
		}
	}
	removeDuplicates(tempAnagramMap)

	// set the first word as a key for the result anagram map
	for key, words := range tempAnagramMap {
		anagramMap[words[0]] = tempAnagramMap[key]
	}

	// Sort the words in each anagram set
	for _, words := range anagramMap {
		sort.Strings(words)
	}
	return anagramMap
}

// sortString sorts the characters of a string and returns the sorted result.
func sortString(s string) string {
	sortedRunes := []rune(s)
	sort.Slice(sortedRunes, func(i, j int) bool {
		return sortedRunes[i] < sortedRunes[j]
	})
	return string(sortedRunes)
}

func removeDuplicates(anagramMap map[string][]string) {
	// Remove sets with only one word
	for key, words := range anagramMap {
		if len(words) == 1 {
			delete(anagramMap, key)
		}
	}
}

// contains checks if a given word is present in the slice of words.
func contains(words []string, word string) bool {
	for _, w := range words {
		if w == word {
			return true
		}
	}
	return false
}
