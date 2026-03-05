package cmd

import (
	"bufio"
	"os"
	"path/filepath"
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

func GetName(fileInfo os.FileInfo) string {

	if !fileInfo.IsDir() {
		return strings.TrimSuffix(filepath.Base(fileInfo.Name()), filepath.Ext(fileInfo.Name()))
	}

	return fileInfo.Name()

}
