package cmd

import "strings"

// Joins strings: "one" "two" "three" to a single string separeated with a string, so "one","two","three"
func ParseArgsIntoSingleString(args []string) string {
	return strings.Join(args, ",")

}

func IsIndexNameValid(index string) bool {
	return !strings.HasPrefix(index, ".") && !strings.Contains(index, "connector")

}

func Reverse[T any](v []T) []T {

	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {

		v[i], v[j] = v[j], v[i]
	}
	return v

}
