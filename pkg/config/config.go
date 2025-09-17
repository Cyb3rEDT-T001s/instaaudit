package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config represents the application configuration
type Config struct {
	Timeout      time.Duration `json:"timeout"`
	MaxWorkers   int           `json:"max_workers"`
	OutputPath   string        `json:"output_path"`
	OutputFormat string        `json:"output_format"`
	Ports        []int         `json:"ports"`
	SkipExploits bool          `json:"skip_exploits"`
	SkipRecon    bool          `json:"skip_recon"`
	Aggressive   bool          `json:"aggressive"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Timeout:      2 * time.Second,
		MaxWorkers:   100,
		OutputPath:   "./audit_report",
		OutputFormat: "json",
		Ports:        []int{}, // Empty means use common ports
		SkipExploits: false,
		SkipRecon:    false,
		Aggressive:   false,
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filepath string) (*Config, error) {
	config := DefaultConfig()

	if filepath == "" {
		return config, nil
	}

	file, err := os.Open(filepath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig saves configuration to a JSON file
func SaveConfig(config *Config, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
}