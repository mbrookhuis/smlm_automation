package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"smlm_automation/pkg/models/inputfile"
)

var AppConfig inputfile.GeneralConfig

// InitConfig reads in inputfile file and ENV variables if set.
func InitConfig() {
	viper.SetConfigName("inputfile")           // name of the input file (without extension)
	viper.SetConfigType("yaml")                // REQUIRED if the inputfile file does not have the extension in the name
	viper.AddConfigPath(".")                   // path to look for the inputfile file in the current directory
	viper.AddConfigPath("$HOME/.your-project") // path to look for the inputfile file in the HOME directory
	viper.AddConfigPath("/etc/your-project/")  // path to look for the inputfile file in /etc/

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using inputfile file:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "No inputfile file found, using defaults or environment variables:", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Unable to unmarshal inputfile:", err)
	}

	// Set default values if not provided in inputfile or env
	if AppConfig.Log.ScreenLevel == "" {
		AppConfig.Log.ScreenLevel = "info" // Default screen level
	}
	if AppConfig.Log.FileLevel == "" {
		AppConfig.Log.FileLevel = "debug" // Default file level
	}
	if AppConfig.Log.FilePath == "" {
		AppConfig.Log.FilePath = "/var/log/do_smlm/do_smlm.log" // Default log file path
	}
}
