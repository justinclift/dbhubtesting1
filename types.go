package main

// Configuration file
type TomlConfig struct {
	Api ApiInfo
}

// Config info for the API server
type ApiInfo struct {
	APIKey string `toml:"api_key"`
	Server string `toml:"server"`
}
