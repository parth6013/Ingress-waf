package utils

import (
	"armor/models"
	"encoding/json"
	"log"
	"os"
)

const (
	CONFIG_PATH = "./config.json"
)

func LoadConfigs() (*models.Config, error) {
	// Read the contents of the file from ./config.json
	data, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Printf("ERROR: Failed to read config file at %s. %s", CONFIG_PATH, err)
		return nil, err
	}

	// Initialize the config with default values
	config := &models.Config{
		Server: models.Server{
			Host: "0.0.0.0",
			Port: 3000,
		},
		Waf: models.Waf{
			EnableCrs:       true,
			CustomRulesPath: "",
			CrsRulesPath:    "./crs4",
		},
		Proxy: models.Proxy{
			TargetHost: "reflector",
			TargetPort: 8080,
		},
		Sam: models.Sam{
			Database:    "telegraf",
			Measurement: "armor",
		},
		Csp: models.Csp{ // Add default CSP values
			Policy: "",
		},
		OIDC: models.OIDC{
			Enable: false,
		},
		Elasticsearch: models.Elasticsearch{
			Enable: false,
			URL:    []string{"http://localhost:9200"},
		},
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Printf("ERROR: Failed to parse %s. %s", CONFIG_PATH, err)
		return nil, err
	}

	log.Printf("INFO: Loaded config: %+v from %s", *config, CONFIG_PATH)

	// Return the Config struct
	return config, nil
}
