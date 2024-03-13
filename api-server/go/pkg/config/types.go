package config

type AppConfig struct {
	Name              string   `json:"appName"`
	Port              int      `json:"port"`
	Prefork           bool     `json:"prefork"`
	Concurrency       int      `json:"concurrency"`
	EnablePrintRoutes bool     `json:"enablePrintRoutes"`
	Debug             bool     `json:"debug"`
	TrustedProxies    []string `json:"trustedProxies"`
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
	ApiBaseUrl   string `json:"apiBaseUrl"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type WordListConfig struct {
	FilePath string `json:"filePath"`
	FileHash string `json:"fileHash"`
}
