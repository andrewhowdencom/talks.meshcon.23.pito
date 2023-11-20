package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/andrewhowdencom/talks.meshcon.23.pito/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrUnableToStartServer = errors.New("unable to start server")
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
	rootCmd.PersistentFlags().Int("go-max-procs", 1, "How many processes to assign the Go runtime")

	// Here, we want to demonstrate how difficult it is to reason through the program absent any sort of logging.
	// Given this, we need it to "swallow" errors in the UI.
	rootCmd.SilenceErrors = true
}

func doPreRun(cmd *cobra.Command, args []string) {
	viper.SetConfigFile(cmd.Flags().Lookup("config").Value.String())

	if err := viper.BindPFlag("go.max-procs", cmd.Flags().Lookup("go-max-procs")); err != nil {
		// Wow it'd be good if we did something here.
	}

	if err := viper.ReadInConfig(); err != nil {
		os.Exit(1)
	}

	// Set runtime constraints
	runtime.GOMAXPROCS(viper.GetInt("go.max-procs"))

}

// doRoot is not yet implemented
func doRoot(cmd *cobra.Command, args []string) error {

	srv := server.New()
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("%w: %s", ErrUnableToStartServer, err)
	}

	return nil
}
