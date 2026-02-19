package test

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type EsFieldsConfig struct {
	Name         string   `yaml:"name"`
	DefaultValue []string `yaml:"default"`
	ValidArgs    []string `yaml:"valid-args"`
	Usage        string   `yaml:"usage"`
}

type EsFields struct {
	Fields []EsFieldsConfig `yaml:"fields"`
}

func TestReadEnvVars(t *testing.T) {
	b, err := os.ReadFile("../cmd/es_fields.yaml")
	if err != nil {
		t.Fatalf("failed to read yaml file: %v", err)
	}

	var e EsFields
	if err := yaml.Unmarshal(b, &e); err != nil {
		t.Fatalf("failed to unmarshal yaml: %v", err)
	}

	if len(e.Fields) == 0 {
		t.Fatal("expected at least one field in yaml config")
	}

	testCmd := &cobra.Command{Use: "testing"}

	for _, f := range e.Fields {
		testCmd.Flags().StringSlice(f.Name, f.DefaultValue, f.Usage)

		err := testCmd.RegisterFlagCompletionFunc(
			f.Name,
			cobra.FixedCompletions(f.ValidArgs, cobra.ShellCompDirectiveNoFileComp),
		)
		if err != nil {
			t.Fatalf("failed to register completion for %s: %v", f.Name, err)
		}
	}

	for _, f := range e.Fields {
		flag := testCmd.Flag(f.Name)
		if flag == nil {
			t.Fatalf("expected flag %s to be registered", f.Name)
		}

		if flag.Usage != f.Usage {
			t.Errorf("flag %s usage mismatch: expected %s, got %s",
				f.Name, f.Usage, flag.Usage)
		}
	}
}

// func filter[T any](s []T, predicate func(T) bool) []T {
// 	result := make([]T, 0, len(s))
// 	for _, v := range s {
// 		if predicate(v) {
// 			result = append(result, v)
// 		}
// 	}
// 	return result
// }
