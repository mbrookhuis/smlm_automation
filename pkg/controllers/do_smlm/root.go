package dosmlm

import (
	"fmt"
	_ "fmt"
	"github.com/spf13/cobra"
	"os"
	model "smlm_automation/pkg/models/inputfile"
	log "smlm_automation/pkg/util/logger"
	ri "smlm_automation/pkg/util/readconfig"
)

var cfgFile string
var AppConfig model.Config

var rootCmd = &cobra.Command{
	Use:   "do_smlm",
	Short: `do_smlm is a CLI to configure SUSE Multi Linux Manager.`,
	Long: `do_smlm is a CLI to configure SUSE Multi Linux Manager. This application uses a configuration file which can be specified with
the --config flag.`,
	SilenceUsage: true,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config before any command runs
		err := initConfig()
		if err != nil {
			return err
		}
		// Initialize logger after config is loaded
		err = log.InitLogger()
		if err != nil {
			return err
		}
		log.Info("Configuration loaded, do-smlm started")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {})
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file (default is config.yaml)")
	cobra.OnFinalize(finalizeRun)
}

func initConfig() error {
	err := ri.ReadConfig(cfgFile, &AppConfig)
	if err != nil {
		return err
	}
	return nil
}

func finalizeRun() {
	// logging.Debugf("PostRun closing logger: do_smlm")
	fmt.Println("PostRun closing logger: do_smlm")
}
