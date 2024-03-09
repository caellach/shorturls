package config

// Config is the configuration for the application
// This is a struct that holds the configuration for the application
// It sets default values for the application, loads from config.json,
// and allows for environment variable overrides

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caellach/shorturl/api-server/go/pkg/env"
)

func DefaultConfigParams() *ConfigParams {
	return &ConfigParams{
		ConfigFilePath: "./config.json",
	}
}

func LoadConfig(configParams *ConfigParams) *Config {
	// Set the default configuration
	config := Config{
		App: AppConfig{
			Name:              "ShortURL Service",
			Port:              3000,
			Prefork:           false,
			Concurrency:       256 * 1024,
			EnablePrintRoutes: true,
			Debug:             false,
		},
		WordList: WordListConfig{
			FilePath: "./wordlist.json",
			FileHash: "REQUIRED",
		},
		MongoDB: MongoDBConfig{
			Uri:      "localhost",
			Port:     27017,
			Database: "shorturls",
			Username: "admin",
			Password: "admin123",
		},
	}

	configFilePath := configParams.ConfigFilePath
	if configFilePath == "" {
		configFilePath = "./config.json"
	}

	// Open the config.json file
	file, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println("Error opening config ("+configFilePath+"):", err)
		return &config
	}
	defer file.Close()

	// Decode the JSON from the file
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding config.json:", err)
		return &config
	}

	// Override the configuration with environment variables
	if os.Getenv("APP_NAME") != "" {
		config.App.Name = os.Getenv("APP_NAME")
	}
	if os.Getenv("APP_PORT") != "" {
		config.App.Port, err = strconv.Atoi(os.Getenv("APP_PORT"))
		if err != nil {
			panic(err)
		}
	}
	if os.Getenv("APP_PREFORK") != "" {
		var prefork = os.Getenv("APP_PREFORK")
		config.App.Prefork = strings.ToLower(prefork) == "true"
	}
	if os.Getenv("APP_CONCURRENCY") != "" {
		config.App.Concurrency, err = strconv.Atoi(os.Getenv("APP_CONCURRENCY"))
		if err != nil {
			panic(err)
		}
		config.App.Concurrency = config.App.Concurrency * 1024
	}
	if os.Getenv("APP_ENABLE_PRINT_ROUTES") != "" {
		var enablePrintRoutes = os.Getenv("APP_ENABLE_PRINT_ROUTES")
		config.App.EnablePrintRoutes = strings.ToLower(enablePrintRoutes) == "true"
	}
	if os.Getenv("APP_DEBUG") != "" {
		var debug = os.Getenv("APP_DEBUG")
		config.App.Debug = strings.ToLower(debug) == "true"
	}

	// Wordlist environment variables
	if os.Getenv("WORDLIST_FILE_PATH") != "" {
		config.WordList.FilePath = os.Getenv("WORDLIST_FILE_PATH")
	}

	// MongoDB environment variables
	if os.Getenv("MONGODB_URI") != "" {
		config.MongoDB.Uri = os.Getenv("MONGODB_URI")
	}
	if os.Getenv("MONGODB_DATABASE") != "" {
		config.MongoDB.Database = os.Getenv("MONGODB_DATABASE")
	}
	if os.Getenv("MONGODB_USERNAME") != "" {
		config.MongoDB.Username = os.Getenv("MONGODB_USERNAME")
	}
	if os.Getenv("MONGODB_PASSWORD") != "" {
		config.MongoDB.Password = os.Getenv("MONGODB_PASSWORD")
	}

	// Discord environment variables
	if os.Getenv("DISCORD_BASE_URL") != "" {
		config.Providers.DiscordConfig.ApiBaseUrl = os.Getenv("DISCORD_API_BASE_URL")
	}
	if os.Getenv("DISCORD_CLIENT_ID") != "" {
		config.Providers.DiscordConfig.ClientID = os.Getenv("DISCORD_CLIENT_ID")
	}
	if os.Getenv("DISCORD_CLIENT_SECRET") != "" {
		config.Providers.DiscordConfig.ClientSecret = os.Getenv("DISCORD_CLIENT_SECRET")
	}

	// Set server env vars
	env.Config.Debug = config.App.Debug

	return &config
}
