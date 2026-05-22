// Package config provides configuration loading and validation for KrillinAI.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the complete application configuration.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	AI       AIConfig       `mapstructure:"ai"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// AIConfig holds AI provider settings.
type AIConfig struct {
	OpenAIKey      string `mapstructure:"openai_key"`
	OpenAIBaseURL  string `mapstructure:"openai_base_url"`
	Model          string `mapstructure:"model"`
	WhisperModel   string `mapstructure:"whisper_model"`
}

// StorageConfig holds file storage settings.
type StorageConfig struct {
	OutputDir string `mapstructure:"output_dir"`
	TempDir   string `mapstructure:"temp_dir"`
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load reads the configuration from the given file path.
// If cfgFile is empty, it searches for a config file in default locations.
func Load(cfgFile string) (*Config, error) {
	v := viper.New()

	setDefaults(v)

	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		v.AddConfigPath(filepath.Join(home, ".krillinai"))
		v.AddConfigPath(".")
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	// Allow environment variable overrides with KRILLIN_ prefix.
	v.SetEnvPrefix("KRILLIN")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found; rely on defaults and environment variables.
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate checks that required configuration values are present.
func (c *Config) Validate() error {
	if c.AI.OpenAIKey == "" {
		return fmt.Errorf("ai.openai_key is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}
	return nil
}

// setDefaults populates sensible default values.
func setDefaults(v *viper.Viper) {
	// Bind to localhost by default instead of all interfaces — safer for local dev.
	v.SetDefault("server.host", "127.0.0.1")
	v.SetDefault("server.port", 8080)
	v.SetDefault("ai.openai_base_url", "https://api.openai.com/v1")
	v.SetDefault("ai.model", "gpt-4o")
	v.SetDefault("ai.whisper_model", "whisper-1")
	v.SetDefault("storage.output_dir", "output")
	v.SetDefault("storage.temp_dir", "temp")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
}
