package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"armor/models"
	"armor/utils"

	"github.com/corazawaf/coraza/v3"
	txhttp "github.com/corazawaf/coraza/v3/http"
	"github.com/corazawaf/coraza/v3/types"
)

var (
	session           = utils.NewSession()
	applicationConfig *models.Config
)

func logError(error types.MatchedRule) {
	log.Printf(
		"%s: %s",
		strings.ToUpper(error.Rule().Severity().String()),
		error.ErrorLog(),
	)

	database, measurementName := applicationConfig.Sam.Database, applicationConfig.Sam.Measurement
	severity, message, data, uri := strings.ToUpper(error.Rule().Severity().String()), error.Message(), error.Data(), error.URI()

	body := []byte(fmt.Sprintf(`%s,severity=%s,influxdb_database=%s message="%s",data="%s",uri="%s"`, measurementName, severity, database, message, data, uri))

	go sendLogs(session, body, nil, nil)
}

func CorazaMiddleware(next http.Handler, config *models.Config) http.Handler {
	log.Printf("INFO: registering CorazaMiddleware with config: %+v", config.Waf)

	session.DissableSSL()
	applicationConfig = config

	wafConfig := coraza.NewWAFConfig().
		WithErrorCallback(logError).
		WithDirectivesFromFile("/etc/crs4/coraza.conf")

	if config.Waf.EnableCrs {
		wafConfig = wafConfig.
			WithDirectivesFromFile("/etc/crs4/crs-setup.conf").
			WithDirectivesFromFile("/etc/crs4/rules/*.conf")
	}

	// TODO: Load additional rules from application configs

	waf, err := coraza.NewWAF(wafConfig)
	if err != nil {
		log.Fatalf("ERROR: Failed to seetup CorazaMiddleware. %s", err)
	}

	return txhttp.WrapHandler(waf, next)
}

func sendLogs(session *utils.Session, body []byte, headers map[string]string, query map[string]string) {
	resp, err := session.Post("https://samv2-svl-fab7-svc-telegraf-regional.cisco.com/write", &body, &headers, &query)
	if err != nil {
		log.Fatalf("Error in sending error log to SAM: %v", err)
		return
	}

	defer resp.Body.Close()
}
