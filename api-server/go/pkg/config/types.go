package config

type AppConfig struct {
	Name              string   `json:"app_name"`
	Port              int      `json:"port"`
	Prefork           bool     `json:"prefork"`
	Concurrency       int      `json:"concurrency"`
	EnablePrintRoutes bool     `json:"enable_print_routes"`
	Debug             bool     `json:"debug"`
	TrustedProxies    []string `json:"trusted_proxies"`
}

type MongoDBConfig struct {
	Uri      string `json:"uri"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenConfig struct {
	Secret string `json:"secret"`
}

// Config is the configuration for the application
type Config struct {
	App       AppConfig       `json:"app"`
	WordList  WordListConfig  `json:"wordlist"`
	MongoDB   MongoDBConfig   `json:"mongodb"`
	Providers ProvidersConfig `json:"providers"`
	Token     TokenConfig     `json:"token"`
}

type ConfigParams struct {
	ConfigFilePath string
}

type ProvidersConfig struct {
	DiscordConfig DiscordConfig `json:"discord"`
}

type DiscordConfig struct {
	ApiBaseUrl   string `json:"api_base_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type WordListConfig struct {
	FilePath string `json:"file_path"`
	FileHash string `json:"file_hash"`
}
