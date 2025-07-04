package inputfile

type Config struct {
	System SystemConfig `mapstructure:"system"`
}

type GeneralConfig struct {
	Log       LogConfig `mapstructure:"log"`
	TimeOut   int       `mapstructure:"timeout"`
	SSL_Check bool      `mapstructure:"ssl_check"`
}

type SystemConfig struct {
	OrganizationName string            `mapstructure:"organization_name"`
	AdminUser        string            `mapstructure:"admin_user"`
	AdminPassword    string            `mapstructure:"admin_password"`
	AdminEmail       string            `mapstructure:"admin_email"`
	SCC              map[string]string `mapstructure:"scc"`
}

type ConfigGeneral struct {
	General GeneralConfig `mapstructure:"general"`
}

type General struct {
	Log       LogConfig `mapstructure:"log"`
	TimeOut   int       `mapstructure:"timeout"`
	SSL_Check bool      `mapstructure:"ssl_check"`
}

type LogConfig struct {
	ScreenLevel string `mapstructure:"screen_level"`
	FileLevel   string `mapstructure:"file_level"`
	FilePath    string `mapstructure:"file_path"`
}

/*
type Config struct {
	General struct {
		Log struct {
			ScreenLevel string `yaml:"screen_level"`
			FileLevel   string `yaml:"file_level"`
			FilePath    string `yaml:"file_path"`
		} `yaml:"log"`
		Timeout  int  `yaml:"timeout"`
		SslCheck bool `yaml:"ssl_check"`
	} `yaml:"general"`
	System struct {
		OrganizationName string `yaml:"organization_name"`
		Scc              struct {
			User1 string `yaml:"user_1"`
			User2 string `yaml:"user_2"`
		} `yaml:"scc"`
		AdminUser     string `yaml:"admin_user"`
		AdminPassword string `yaml:"admin_password"`
		AdminEmail    string `yaml:"admin_email"`
	} `yaml:"system"`
}
*/

/*
general:
  log:
    screen_level: info
    file_level: debug
    file_path: /var/log/do_smlm/do_smlm.log
  timeout: 1200
  ssl_certificate_check: True

system:
  organization_name: orgname
  scc:
    user_1: password_1
    user_2: password_2
  admin_user: sm-admin
  admin_password: SUSE4ever!






*/
