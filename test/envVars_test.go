package test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"sigs.k8s.io/yaml"
)

type EsFIELDSConfig struct {
	Name        string   `json:"name"`
	Validargs   []string `json:"validargs"`
	Description string   `json:"description"`
}

type EsFields struct {
	Fields []EsFIELDSConfig `json:"fields"`
}

func TestReadEnvVars(t *testing.T) {

	b, err := os.ReadFile("../es_fields.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var e EsFields
	err = yaml.Unmarshal(b, &e)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", e)

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
