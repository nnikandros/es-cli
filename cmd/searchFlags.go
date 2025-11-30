package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

type SearchFlags struct {
	Size     int
	Fields   []string
	Tabular  bool
	Time     bool
	Reverse  bool
	Terms    bool
	Id       []string
	LEVEL    []string
	APP_NAME []string
}

// var validLevelStates = map[string]struct{}{"DEBUG": struct{}{}, "ERROR": struct{}{}, "INFO": struct{}{}}
var validLevelStates = []string{"DEBUG", "ERROR", "INFO"}

func ParsedFlagsFromCmd(cmd *cobra.Command) (SearchFlags, error) {

	size, _ := cmd.Flags().GetInt("size")

	fields, _ := cmd.Flags().GetStringSlice("fields")

	resizeSize := min(size, 10000)

	time, _ := cmd.Flags().GetBool("time")
	reverse, _ := cmd.Flags().GetBool("reverse")
	tabular, _ := cmd.Flags().GetBool("tab")

	terms, _ := cmd.Flags().GetBool("terms")
	id, _ := cmd.Flags().GetStringSlice("id")
	level, _ := cmd.Flags().GetStringSlice("LEVEL")
	app_name, _ := cmd.Flags().GetStringSlice("APP_NAME")

	for _, l := range level {
		ok := slices.Contains(validLevelStates, l)
		if !ok {
			return SearchFlags{}, fmt.Errorf("level provided %v is not supported. levels supported %v", level, validLevelStates)
		}
	}

	if !time && reverse {
		return SearchFlags{}, fmt.Errorf("you have provided revese but not time")
	}

	if len(id) > 0 && !terms {
		return SearchFlags{}, fmt.Errorf("you have provided a field but not a terms")

	}

	return SearchFlags{Size: resizeSize, Fields: fields, Time: time, Reverse: reverse, Tabular: tabular, Terms: terms, Id: id, LEVEL: level, APP_NAME: app_name}, nil

}

// types.FieldAndFormat
