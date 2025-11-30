package cmd

import "strings"

// Joins strings: "one" "two" "three" to a single string separeated with a string, so "one","two","three"
func ParseArgsIntoSingleString(args []string) string {
	return strings.Join(args, ",")

}

func IsIndexNameValid(index string) bool {
	return !strings.HasPrefix(index, ".") && !strings.Contains(index, "connector")

}
