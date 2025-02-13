package models

type Server struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type Waf struct {
	EnableCrs       bool   `json:"enableCrs,omitempty"`
	CustomRulesPath string `json:"customRulesPath,omitempty"`
}

type Proxy struct {
	TargetHost string `json:"targetHost,omitempty"`
	TargetPort int    `json:"targetPort,omitempty"`
}

type Sam struct {
	Database    string `json:"database,omitempty"`
	Measurement string `json:"measurement,omitempty"`
}

type Csp struct {
	Policy string `json:"policy,omitempty"` // Define the CSP policy as a string
}

type OIDC struct {
	Enable              bool   `json:"enable,omitempty"`
	ClientID            string `json:"clientID,omitempty"`
	ProviderUrl         string `json:"providerUrl,omitempty"`
	RedirectApplication string `json:"redirectApplication,omitempty"`
}

// Config holds the entire application configuration
type Config struct {
	Server Server `json:"server,omitempty"`
	Waf    Waf    `json:"waf,omitempty"`
	Proxy  Proxy  `json:"proxy,omitempty"`
	Sam    Sam    `json:"sam,omitempty"`
	Csp    Csp    `json:"csp,omitempty"` // Add the Csp struct here
	OIDC   OIDC   `json:"oidc,omitempty"`
}
