package dosmlm

import (
	"github.com/spf13/cobra"
	log "smlm_automation/pkg/util/logger"
	// logging "smlm_automation/pkg/util/logger"
)

var doSystemCmd = &cobra.Command{
	Use:   "do-system",
	Short: "do-system set SMLM general configuration",
	Long:  `do-system set SMLM general configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeDoSystem()
	},
}

func init() {
	rootCmd.AddCommand(doSystemCmd)
}

func executeDoSystem() error {
	log.Debug("doSystem called")
	log.Infof("Password  :  %s\n", AppConfig.System.AdminPassword)

	return nil
}
