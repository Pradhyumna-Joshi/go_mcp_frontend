package models

type ModelConfig struct {
	Name     string
	Provider string
	Model    string
	APIKey   string
	BaseURL  string
	HasSaved bool
}

type MCPServerConfig struct {
	Name      string
	BaseURL   string
	Transport string
	HasSaved  bool
}
