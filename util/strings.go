package util

import "strings"

func Capitalize(a string) string {
	splitter := strings.Split(a, " ")
	for i, v := range splitter {
		splitter[i] = mapUpper(v)
	}
	return strings.Join(splitter, " ")
}

func mapUpper(s string) string {
	first, rest := string(s[0]), s[1:]
	return strings.ToUpper(first) + rest
}
