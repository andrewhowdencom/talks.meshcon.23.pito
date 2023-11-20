package cmd

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"log/slog"

	"github.com/andrewhowdencom/talks.meshcon.23.pito/server"
	"github.com/andrewhowdencom/talks.meshcon.23.pito/telemetry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.opentelemetry.io/otel"
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

func init() {
	// Define the path for the configuration file.
	rootCmd.PersistentFlags().String("config", "/etc/.talks.meshcon.23.pito.yaml",
		"config file (default is /etc/.talks.meshcon.23.pito.yaml)")
	rootCmd.PersistentFlags().String("listen-address", "localhost:80",
		"The address on which the server should listen")
	rootCmd.PersistentFlags().Int("go-max-procs", 1, "How many processes to assign the Go runtime")

	// Here, we want to demonstrate how difficult it is to reason through the program absent any sort of logging.
	// Given this, we need it to "swallow" errors in the UI.
	rootCmd.SilenceErrors = true
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("failed to start application", "error", err)
		os.Exit(1)
	}
}

// doPreRun is a series of persistent pre-execution thing that should happen. Could also go in init().
func doPreRun(cmd *cobra.Command, args []string) {
	viper.SetConfigFile(cmd.Flags().Lookup("config").Value.String())

	if err := viper.BindPFlag("go.max-procs", cmd.Flags().Lookup("go-max-procs")); err != nil {
		slog.Info("failed to bind flag", "flag", "go-max-procs", "error", err)
	}

	if err := viper.BindPFlag("server.listen-address", cmd.Flags().Lookup("listen-address")); err != nil {
		slog.Info("failed to bind flag", "flag", "server.listen-address", "error", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("failed read configuration", "error", err)
		os.Exit(1)
	}

	// Setup a global tracer provider, based on this applications configuration.
	tp, err := telemetry.NewTracerProvider()
	if err != nil {
		slog.Error("failed to setup tracing", "error", err)
	}

	otel.SetTracerProvider(tp)

	// Set runtime constraints
	runtime.GOMAXPROCS(viper.GetInt("go.max-procs"))

}

// doRoot starts the server
func doRoot(cmd *cobra.Command, args []string) error {

	srv := server.New(server.WithListenAddr(viper.GetString("server.listen-address")))
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("%w: %s", ErrUnableToStartServer, err)
	}

	return nil
}
