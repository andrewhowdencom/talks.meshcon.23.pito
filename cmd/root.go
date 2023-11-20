package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:    "rabbit",
	Short:  "An application that, when queried over HTTP, returns a text Rabbit",
	PreRun: doPreRun,
	RunE:   doRoot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Define the path for the configuration file.
	rootCmd.PersistentFlags().String("config", "/etc/.talks.meshcon.23.pito.yaml", "config file (default is /etc/.talks.meshcon.23.pito.yaml)")

	// Here, we want to demonstrate how difficult it is to reason through the program absent any sort of logging.
	// Given this, we need it to "swallow" errors in the UI.
	rootCmd.SilenceErrors = true
}

func doPreRun(cmd *cobra.Command, args []string) {
	viper.SetConfigFile(cmd.Flags().Lookup("config").Value.String())
	if err := viper.ReadInConfig(); err != nil {
		os.Exit(1)
	}
}

// doRoot is not yet implemented
func doRoot(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("not implemented")
}
