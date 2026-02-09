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
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
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
es index create --directory ./elasticops/settings_mappings/test-index			
`,
	}

	return cmd

}

func addCreateFlags(create *cobra.Command) *cobra.Command {
	create.Flags().StringP("directory", "d", "", "path to the directory where you have the mappings and settings")
	create.Flags().StringP("path", "p", "", "path to the json file where you have the mappings and settings together")
	return create

}

func runCreateIndexCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		mappings := types.TypeMapping{}
		settings := types.IndexSettings{}

		directory, _ := cmd.Flags().GetString("directory")
		pathToJson, _ := cmd.Flags().GetString("path")

		var indexName string

		switch len(args) {
		case 0:
			indexName = filepath.Base(directory)
		case 1:
			indexName = args[0]
		default:
			return fmt.Errorf("too many arguments")

		}

		if len(directory) != 0 {
			absPathDirToMappingsAndSettings, _ := filepath.Abs(directory)

			_, err := os.Stat(absPathDirToMappingsAndSettings)
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return fmt.Errorf("the path %v does not exist, %w", directory, err)
				}
				return fmt.Errorf("stat of the %v returned an  error %w", directory, err)
			}

			pathToMappings := filepath.Join(absPathDirToMappingsAndSettings, "mappings.json")
			pathToSettings := filepath.Join(absPathDirToMappingsAndSettings, "settings.json")

			// to use serde.DecodeV2
			mappings, err1 := serde.DecodeJsonFileToStruct[types.TypeMapping](pathToMappings)
			settings, err2 := serde.DecodeJsonFileToStruct[types.IndexSettings](pathToSettings)

			if err1 != nil || err2 != nil {
				return fmt.Errorf("parsing the mappings or settings json files %w, %w", err1, err2)
			}

		}
		r, err := createIndex(es, indexName, mappings, settings)
		if err != nil {
			return fmt.Errorf("at creating index: %w", err)
		}

		if err := json.NewEncoder(cmd.OutOrStderr()).Encode(r); err != nil {
			return serde.SerializingError(err)
		}

		return nil
	}

}

func createIndex(es *elasticsearch.TypedClient, indexName string, mappings types.TypeMapping, settings types.IndexSettings) (*create.Response, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	r, err := es.Indices.Create(indexName).Mappings(&mappings).Settings(&settings).Do(ctx)
	if err != nil {
		return r, fmt.Errorf("at creating index %w", err)
	}

	return r, nil
}

type MappingsAndSettings struct {
	Mappings types.TypeMapping   `json:"mappings"`
	Settings types.IndexSettings `json:"settings"`
}

func handleDirCase(directory string) (MappingsAndSettings, error) {
	absPathDirToMappingsAndSettings, _ := filepath.Abs(directory)

	_, err := os.Stat(absPathDirToMappingsAndSettings)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return MappingsAndSettings{}, fmt.Errorf("the path %v does not exist, %w", directory, err)
		}
		return MappingsAndSettings{}, fmt.Errorf("stat of the %v returned an  error %w", directory, err)
	}

	pathToMappings := filepath.Join(absPathDirToMappingsAndSettings, "mappings.json")
	pathToSettings := filepath.Join(absPathDirToMappingsAndSettings, "settings.json")

	// to use serde.DecodeV2
	mappings, err1 := serde.DecodeJsonFileToStruct[types.TypeMapping](pathToMappings)
	settings, err2 := serde.DecodeJsonFileToStruct[types.IndexSettings](pathToSettings)

	if err1 != nil || err2 != nil {
		return MappingsAndSettings{}, fmt.Errorf("parsing the mappings or settings json files %w, %w", err1, err2)
	}

	return MappingsAndSettings{mappings, settings}, nil
}

func handleFileCase(file string) (MappingsAndSettings, error) {
	absPathDirToMappingsAndSettings, _ := filepath.Abs(file)

	_, err := os.Stat(absPathDirToMappingsAndSettings)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return MappingsAndSettings{}, fmt.Errorf("the path %v does not exist, %w", file, err)
		}
		return MappingsAndSettings{}, fmt.Errorf("stat of the %v returned an  error %w", file, err)
	}

}
