package dosmlm

import (
	"github.com/spf13/cobra"
	logging "smlm_automation/pkg/util/logger"
)

var doSystemCmd = &cobra.Command{
	Use:   "doSystemCmd",
	Short: "doSystemCmd set SMLM general configuration",
	Long:  `doSystemCmd set SMLM general configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeDoSystem()
	},
}

func init() {
	rootCmd.AddCommand(doSystemCmd)
}

func executeDoSystem() error {
	logging.Debug("doSystem called")

	return nil
}
