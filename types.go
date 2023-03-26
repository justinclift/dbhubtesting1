package main

// TomlConfig is the overall container structure for our config file
type TomlConfig struct {
	Api ApiInfo
}

// ApiInfo holds configuration info for the API server
type ApiInfo struct {
	APIKey     string `toml:"api_key"`
	Server     string `toml:"server"`
	VerifyCert bool   `toml:"verify_cert"`
}
