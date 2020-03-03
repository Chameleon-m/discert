package apiserver

type Config struct {
	BindAddr       string `toml:"bind_addr"`
	LogLevel       string `toml:"log_level"`
	DatabaseDriver string `toml:"database_driver"`
	DatabaseURL    string `toml:"database_url"`
	SessionKey     string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}
