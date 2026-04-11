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

type Input struct {
	Id          string
	Name        string
	Label       string
	Value       string
	Type        string
	Title       string
	Required    bool
	ReadOnly    bool
	PlaceHolder string
	ClassName   string
}

type Button struct {
	Id        string
	Type      string
	ClassName string
	Content   string
}

type NavItem struct {
	GetURL    string
	Target    string
	Swap      string
	IconClass string
	Title     string
}
