package cmd

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

type MigrateCmd = *cobra.Command

func migrateCmdFunc(es *elasticsearch.TypedClient) SearchCmd {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrations for es",
		Long:  "Running a search query against an index\nThe arguments can be an index a list of indexes separated by space or index name with wildcard. Can also be _all to search all indices",
		Args:  cobra.NoArgs,
	}

	return cmd
}

func migrateFunc(es *elasticsearch.TypedClient) MigrateCmd {

	migrateCmd := migrateCmdFunc(es)

	cloneSubCommand := addCloneFlags(cloneMigrateCmdFunc(es))
	reindexSubComand := addReindexFlags(reindexMigrateCmdFunc(es))
	migrateCmd.AddCommand(cloneSubCommand)
	migrateCmd.AddCommand(reindexSubComand)

	return migrateCmd
}
