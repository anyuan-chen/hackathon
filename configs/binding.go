package configs

type Config struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	DatabaseUrl string `mapstructure:"databaseurl"`
}
