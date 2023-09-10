package apiserver

type Config struct {
	BindAddr  string `toml:"bind_addr"`
	SecretKey string `toml:"secret_key"`
	ApiKey    string `toml:"api_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}
