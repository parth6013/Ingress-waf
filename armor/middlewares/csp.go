package middlewares

import (
	"armor/models"
	"log"
	"net/http"

	"github.com/unrolled/secure"
)

// CSPMiddleware creates a middleware to apply the CSP policy.
func CSPMiddleware(next http.Handler, config *models.Config) http.Handler {
	// Log the CSP configuration being used
	log.Printf("INFO: registering CSPMiddleware with config: %+v", config.Csp)

	// Create a new instance of the secure middleware with the constructed CSP
	secureMiddleware := secure.New(secure.Options{
		ContentSecurityPolicy: config.Csp.Policy, // Use the CSP string from the config
	})

	return secureMiddleware.Handler(next)
}
