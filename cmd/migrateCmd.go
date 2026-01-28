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
		Long:  "Subcommand to help with migrating indices.",
		Args:  cobra.NoArgs,
	}

	return cmd
}

func migrateFunc(es *elasticsearch.TypedClient) MigrateCmd {

	migrateCmd := migrateCmdFunc(es)

	cloneSubCommand := addCloneFlags(cloneMigrateCmdFunc(es), es)
	reindexSubComand := addReindexFlags(reindexMigrateCmdFunc(es), es)
	migrateCmd.AddCommand(cloneSubCommand)
	migrateCmd.AddCommand(reindexSubComand)

	return migrateCmd
}
