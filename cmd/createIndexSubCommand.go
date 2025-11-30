package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"serde"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/spf13/cobra"
)

func createIndexCmdFunc(es *elasticsearch.TypedClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create subcommand creates an index",
		Long: `create subcommand creates an index with args the name of the new index and as flag
argument the directory where you store the mappings.json and settings.json`,
		RunE: runCreateIndexCmdFunc(es),
		Example: `es index create <name-new-index> -d ./elasticops/settings_mappings/test-index
es index create <name-new-index> --directory ./elasticops/settings_mappings/test-index		
`,
	}

	return cmd

}

func addCreateFlags(create *cobra.Command) *cobra.Command {
	create.Flags().StringP("directory", "d", "", "path to the directory where you have the mappings and settings")
	return create

}

func runCreateIndexCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return fmt.Errorf("the subcommand create takes exactly one argument")
		}

		indexName := args[0]

		directory, err := cmd.Flags().GetString("directory")
		if err != nil {
			return fmt.Errorf("at reading the value of the path flag %w", err)
		}

		absPathDirToMappingsAndSettings, _ := filepath.Abs(directory)

		_, err = os.Stat(absPathDirToMappingsAndSettings)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("the path %v does not exist, %w", directory, err)
			}
			return fmt.Errorf("stat of the %v returned an  error %w", directory, err)
		}

		pathToMappings := filepath.Join(absPathDirToMappingsAndSettings, "mappings.json")
		pathToSettings := filepath.Join(absPathDirToMappingsAndSettings, "settings.json")

		mappings, err1 := serde.DecodeJsonFileToStruct[types.TypeMapping](pathToMappings)
		settings, err2 := serde.DecodeJsonFileToStruct[types.IndexSettings](pathToSettings)

		if err1 != nil || err2 != nil {
			return fmt.Errorf("parsing the mappings or settings json files %w, %w", err1, err2)
		}

		return createIndex(es, indexName, mappings, settings)
	}

}

func createIndex(es *elasticsearch.TypedClient, indexName string, mappings types.TypeMapping, settings types.IndexSettings) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := es.Indices.Create(indexName).Mappings(&mappings).Settings(&settings).Do(ctx)
	if err != nil {
		return fmt.Errorf("at creating index %w", err)
	}

	marshaled, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("at serializing the response %w", err)
	}
	fmt.Printf("%s\n", marshaled)

	return nil
}

// func MultipeErrorHandler(errors ...int) error {
// 	fmt.Println(nums)

// 	return nil

// }
