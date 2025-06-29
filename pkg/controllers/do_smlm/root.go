package dosmlm

import (
	_ "fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"smlm_automation/pkg/config"

	logging "smlm_automation/pkg/util/logger"
	returncodes "smlm_automation/pkg/util/returnCodes"
)

var rootCmd = &cobra.Command{
	Use:   "do_smlm",
	Short: `do_smlm is a CLI to configure SUSE Multi Linux Manager.`,
	Long: `do_smlm is a CLI to configure SUSE Multi Linux Manager. This application uses a configuration file which can be specified with
the --config flag.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config before any command runs
		config.InitConfig()
		// Initialize logger after config is loaded
		logging.InitLogger()
		return nil
	},
}
var cfgFile string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logging.Fatalf("%s: %s\n%s", returncodes.ErrRunningService, "do_smlm", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file (default is config.yaml)")
	err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	if err != nil {
		return
	}
	cobra.OnFinalize(finalizeRun)
}

func initConfig() {
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		logging.Debugf("Custom logger initialized: do_smlm")
	}
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err) // cobra.CheckErr is a helper that calls os.Exit(1) on error

		// Search config in home directory with name ".yourcli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".yourcli") // Searches for .yourcli.yaml, .yourcli.json, etc.
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logging.Info(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		// Handle cases where the config file is not found.
		// This is often not an error, as you might rely on flags or env vars.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file isn't found; ignore error
			logging.Warn(os.Stderr, "No config file found. Using default values or environment variables.")
		} else {
			// Some other error occurred when reading the config file
			logging.Error(os.Stderr, "Error reading config file: %s\n", err)
			os.Exit(1)
		}
	}
}

// PostRun functions seem not to run reliably at least when I tested
// See: https://github.com/spf13/cobra/issues/914
func finalizeRun() {
	logging.Debugf("PostRun closing logger: do_smlm")
}
