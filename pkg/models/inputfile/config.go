package inputfile

type LogConfig struct {
	ScreenLevel string `mapstructure:"screen_level"`
	FileLevel   string `mapstructure:"file_level"`
	FilePath    string `mapstructure:"file_path"`
}

type GeneralConfig struct {
	Log LogConfig `mapstructure:"log"`
}

type Config struct {
	General GeneralConfig `mapstructure:"general"`
	// Add other configuration parameters here
}
