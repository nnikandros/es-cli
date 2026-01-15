package test

import (
	"fmt"
	"log"
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
		log.Fatal(err)
	}

	var e EsFields
	err = yaml.Unmarshal(b, &e)
	if err != nil {
		log.Fatal(err)
	}

	testCmd := &cobra.Command{Use: "testing additiong of flags"}

	for _, f := range e.Fields {
		testCmd.Flags().StringSlice(f.Name, f.DefaultValue, f.Usage)
		testCmd.RegisterFlagCompletionFunc(f.Name, cobra.FixedCompletions(f.ValidArgs, cobra.ShellCompDirectiveNoFileComp))
		fmt.Printf("%+v", testCmd.Flag(f.Name))
	}

	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	t.Error(err)
	// }

	// // for _, e := range os.Environ() {
	// // 	// pair := strings.SplitN(e, "=", 2)
	// // 	// fmt.Println(pair[0])
	// // 	fmt.Println(e)
	// // }

	// s := filter(os.Environ(), func(x string) bool {
	// 	return strings.Contains(x, "ES_FIELD")
	// })

	// fmt.Println(strings.Split(s[0], "="))
}

func filter[T any](s []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(s)) // Pre-allocate for efficiency
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
