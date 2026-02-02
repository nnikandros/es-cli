/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "es",
	Short: "Elasticsearch Operations",
	Long: `The es cli brings to you some of Elasticsearch operations in your terminal.
Make sure you have the following environment variables exported
- ES_API_KEY: not much to explain, create a token from Kibana
- ES_HOSTS: (can be one host or a comma separated, for example ES_HOSTS=HOST-1,HOST-2,HOST-3 )
- ES_CERT_PATH: Path to your elasticsearch_ca.crt to make the connection secure. For example ~/<user>/elasticsearch_ca.crt
Certain operations for example pinging the cluster and listing all the indices require certain privileges that your API key may or may not have.
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// err := godotenv.Load()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.escobra.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	typedClient, err := NewElasticTypedClient()
	if err != nil {
		log.Fatal(err)
	}

	indicesCmd := IndexCmdFunc(typedClient)

	searchCmd := addSearchFlags(searchCmdFunc(typedClient))
	// countcmd := addCountFlags(countCmdFunc(typedClient))

	clusterCmd := InfoCmdFunc(typedClient)

	migrateCmd := migrateFunc(typedClient)

	RootCmd.AddCommand(indicesCmd)
	RootCmd.AddCommand(searchCmd)
	RootCmd.AddCommand(clusterCmd)
	// RootCmd.AddCommand(countcmd)
	RootCmd.AddCommand(migrateCmd)

	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
