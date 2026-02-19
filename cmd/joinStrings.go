package cmd

import (
	"bufio"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func IsIndexNameValid(index string) bool {
	return !strings.HasPrefix(index, ".") && !strings.Contains(index, "connector")

}

func Reverse[T any](v []T) []T {

	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {

		v[i], v[j] = v[j], v[i]
	}
	return v

}

// // Parse args into a single string. Usually (Always?) args for the es cli are index names. If we provide args then we simply join them. If we don't provide then we read them from  Stdin and  join them with a comma.
// func ParseArgsIntoString(cmd *cobra.Command, args []string) (index string) {
// 	indexList := make([]string, 0)

// 	if len(args) > 0 {
// 		index = strings.Join(args, ",")
// 	} else {
// 		scanner := bufio.NewScanner(cmd.InOrStdin())
// 		for scanner.Scan() {
// 			indexList = append(indexList, scanner.Text())
// 		}
// 		index = strings.Join(indexList, ",")
// 	}

// 	return
// }

func ParseArgsIntoString(cmd *cobra.Command, args []string) string {
	if len(args) > 0 {
		return strings.Join(args, ",")
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return ""
	}

	var indexList []string
	scanner := bufio.NewScanner(cmd.InOrStdin())
	for scanner.Scan() {
		indexList = append(indexList, scanner.Text())
	}

	return strings.Join(indexList, ",")
}
