package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"serde"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

type IndexCmd = *cobra.Command

type RunEFunc func(cmd *cobra.Command, args []string) error

func indicesCmdFunc(es *elasticsearch.TypedClient) IndexCmd {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "The index subcommands offers various operations related to indices.",
		Long:  `A longer description that spans multiple lines and likely contains`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		// Run: runIndices,
		RunE: runIndicesCmdFunc(es),
		Example: `es index <index-name> --ping
es index <index-name-v*> --mappings
es index <test-log-v4,test-index-v5> --mappings
es index <test-log-v4, test-log-v5> --mappings
es index <test-log-v*> --settings
es index <test-log-v*> -sm
`,
		SilenceUsage:      true,
		ValidArgsFunction: ValidArgsFuncAutoCompletion(es),
	}

	return cmd

}

func addIndicesFlags(indicesCmd IndexCmd) IndexCmd {
	indicesCmd.Flags().BoolP("mappings", "m", false, "mappings flag. Enable it will print the mappings to stdout")
	indicesCmd.Flags().BoolP("settings", "s", false, "settings flag. Enable it will print the settings to stdout")
	indicesCmd.Flags().BoolP("ping", "p", false, "ping flag. Pining the index asserts that you can connect to the index. It makes a match_all query and assets that the response is OK")
	return indicesCmd

}

// Function wrapper to create the index cmd. Adds the flags for the index cmd, all of its subcommands and all their flags too.
func IndexCmdFunc(es *elasticsearch.TypedClient) IndexCmd {
	indicesCmd := addIndicesFlags(indicesCmdFunc(es))

	listSubcommand := addListFlags(listIndicesCmdFunc(es))
	createSubCommand := addCreateFlags(createIndexCmdFunc(es))
	deleteSubCommand := deleteIndexCmdFunc(es)
	cloneSubCommand := addCloneFlags(cloneIndexCmdFunc(es))

	indicesCmd.AddCommand(listSubcommand)
	indicesCmd.AddCommand(createSubCommand)
	indicesCmd.AddCommand(deleteSubCommand)
	indicesCmd.AddCommand(cloneSubCommand)

	return indicesCmd

}

func runIndicesCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		// if len(args) == 0 {
		// 	return fmt.Errorf("please provide at least one index name")
		// }
		var index string
		if len(args) > 0 {
			index = ParseArgsIntoSingleString(args)
		} else {
			scanner := bufio.NewScanner(cmd.InOrStdin())
			if scanner.Scan() {
				index = scanner.Text()
			}
		}

		mappings, _ := cmd.Flags().GetBool("mappings")
		settings, _ := cmd.Flags().GetBool("settings")
		ping, _ := cmd.Flags().GetBool("ping")

		switch {
		case mappings && settings:
			b, err := getIndexInfo(es, index)
			if err != nil {
				return fmt.Errorf("at getting index info %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s", b)

			return nil

		case mappings:
			b, err := getIndexMapping(es, index)
			if err != nil {
				return fmt.Errorf("at getting the index mapping %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s", b)
			return nil
		case settings:
			b, err := getIndexSettings(es, index)
			if err != nil {
				return fmt.Errorf("at getting the index settings %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s", b)
			return nil
		case ping:
			err := pingIndex(es, index)
			if err != nil {
				return fmt.Errorf("at ping index %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "OK\n")
			return nil

		default:
			b, err := getIndexInfo(es, index)
			if err != nil {
				return fmt.Errorf("at getting index info %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "%s", b)
			return nil

		}

	}
}

// https://stackoverflow.com/questions/67214264/how-do-i-pass-the-database-connection-to-all-cobra-commands

func getIndexInfo(es *elasticsearch.TypedClient, indexName string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := es.Indices.Get(indexName).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("at doing a get at the index %s with error %w", indexName, err)
	}

	b, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return nil, serde.SerializingError(err)
	}
	return b, nil

}

// indexName supports wildcards (`*`). To target all data streams and indices, omit this parameter or use `*` or `_all`. API Name: index
func getIndexMapping(es *elasticsearch.TypedClient, indexName string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := es.Indices.GetMapping().Index(indexName).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("at getting the mapping of index: %s with error %w", indexName, err)
	}

	b, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return nil, fmt.Errorf("at serializing response %w", err)
	}

	return b, nil

}

// can be comma separated string
func getIndexSettings(es *elasticsearch.TypedClient, indexName string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r, err := es.Indices.GetSettings().Index(indexName).Do(ctx)
	if err != nil {
		return nil, err
	}

	b, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return nil, err
	}

	return b, nil

}

func pingIndex(es *elasticsearch.TypedClient, indexName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := es.Search().Index(indexName).Size(1).Do(ctx)
	if err != nil {
		return fmt.Errorf("pinging the index %s returned: %w", indexName, err)
	}

	return nil
}

// indicesCmd := addIndicesFlags(indicesCmdFunc(typedClient))
// 	listSubcommand := addListFlags(listIndicesCmdFunc(typedClient))
// 	createSubCommand := addCreateFlags(createIndexCmdFunc(typedClient))
// 	deleteSubCommand := deleteIndexCmdFunc(typedClient)

// 	indicesCmd.AddCommand(listSubcommand)
// 	indicesCmd.AddCommand(createSubCommand)
// 	indicesCmd.AddCommand(deleteSubCommand)
