package cmd

import (
	"bufio"
	"strings"

	"github.com/spf13/cobra"
)

func IsIndexNameValid(index string) bool {
	return !strings.HasPrefix(index, ".") && !strings.Contains(index, "connector")

}

func ParseArgsIntoString(cmd *cobra.Command, args []string) string {
	var index string
	// var indexList []string
	indexList := make([]string, 0)

	if len(args) > 0 {
		index = strings.Join(args, ",")
	} else {
		scanner := bufio.NewScanner(cmd.InOrStdin())
		for scanner.Scan() {
			indexList = append(indexList, scanner.Text())
		}
		index = strings.Join(indexList, ",")
	}

	return index
}
