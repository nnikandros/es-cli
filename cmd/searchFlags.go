package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type SearchFlags struct {
	Size           int
	Fields         []string
	Tabular        bool
	Time           bool
	Reverse        bool
	Terms          bool
	Id             []string
	FieldsTermsMap []EsFieldsConfig
	Should         bool
	Must           bool
	Not            bool
}

func ParsedFlagsFromCmd(cmd *cobra.Command) (SearchFlags, error) {

	size, _ := cmd.Flags().GetInt("size")

	fields, _ := cmd.Flags().GetStringSlice("fields")

	resizeSize := min(size, 10000)

	time, _ := cmd.Flags().GetBool("time")
	reverse, _ := cmd.Flags().GetBool("reverse")
	tabular, _ := cmd.Flags().GetBool("tab")

	terms, _ := cmd.Flags().GetBool("terms")
	id, _ := cmd.Flags().GetStringSlice("id")

	o := make([]EsFieldsConfig, 0, len(e.Fields))
	for _, f := range e.Fields {
		valuesForField, _ := cmd.Flags().GetStringSlice(f.Name)
		o = append(o, EsFieldsConfig{Name: f.Name, Value: valuesForField})
	}

	// if len(o) > 1 {
	// 	return SearchFlags{}, fmt.Errorf("you have provided more than one term. Currently one is supported")
	// }

	if !time && reverse {
		return SearchFlags{}, fmt.Errorf("you have provided revese but not time")
	}

	return SearchFlags{Size: resizeSize, Fields: fields, Time: time, Reverse: reverse, Tabular: tabular, Terms: terms, Id: id, FieldsTermsMap: o}, nil

}
