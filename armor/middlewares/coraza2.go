package middlewares

import (
	"log"
	"net/http"
	"strings"
	"time"

	"armor/models"

	"github.com/corazawaf/coraza/v3"
	txhttp "github.com/corazawaf/coraza/v3/http"
	"github.com/corazawaf/coraza/v3/types"
)

var (
	applicationConfig *models.Config
)

// LogEntry represents the structure for Elasticsearch logging

func logError(error types.MatchedRule) {
	severity := strings.ToUpper(error.Rule().Severity().String())

	log.Printf(
		"%s: %s",
		severity,
		error.ErrorLog(),
	)

	message, _, _ := error.Message(), error.Data(), error.URI()

	// Create log entry for Elasticsearch
	logEntry := LogEntryCoraza{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Severity:  severity,
		Message:   message,
	}

	// Send to Elasticsearch asynchronously
	go SendToElasticSearchCoraza(logEntry)
}

func CorazaMiddleware(next http.Handler, config *models.Config) http.Handler {
	log.Printf("INFO: registering CorazaMiddleware with config: %+v", config.Waf)

	applicationConfig = config

	wafConfig := coraza.NewWAFConfig().
		WithErrorCallback(logError).
		WithDirectivesFromFile("/crs4/coraza.conf")

	if config.Waf.EnableCrs {
		wafConfig = wafConfig.
			WithDirectivesFromFile("/crs4/crs-setup.conf").
			WithDirectivesFromFile("/crs4/rules/*.conf")
	}

	// TODO: Load additional rules from application configs

	waf, err := coraza.NewWAF(wafConfig)
	if err != nil {
		log.Fatalf("ERROR: Failed to setup CorazaMiddleware. %s", err)
	}

	return txhttp.WrapHandler(waf, next)
}
