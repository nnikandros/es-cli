package cmd

import (
	"bufio"
	"strings"

	"github.com/spf13/cobra"
)

func IsIndexNameValid(index string) bool {
	return !strings.HasPrefix(index, ".") && !strings.Contains(index, "connector")

}

// Parse args into a single string. Usually (Always?) args for the es cli are index names. If we provide args then we simply join them. If we don't provide then we read them from  Stdin and  join them with a comma.
func ParseArgsIntoString(cmd *cobra.Command, args []string) (index string) {
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

	return
}
