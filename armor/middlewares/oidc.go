package middlewares

import (
	"armor/models"
	"armor/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	oidcNonceLength = 16
	oidcStateLength = 16
	callbackUrl     = "/auth/oidc/callback"
)

var (
	clientSecret = os.Getenv("DUO_OAUTH2_CLIENT_SECRET")
)

func loginRedirect(w http.ResponseWriter, r *http.Request, oauth2Config *oauth2.Config) {
	state, err := utils.RandString(oidcStateLength)
	if err != nil {
		log.Printf("ERROR: could not generate random string for oidc state. %s", err)
		utils.ErrorJsonResponse(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	nonce, err := utils.RandString(oidcNonceLength)
	if err != nil {
		log.Printf("ERROR: could not generate random string for oidc nonce. %s", err)
		utils.ErrorJsonResponse(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	utils.SetCookie(w, r, "state", state)
	utils.SetCookie(w, r, "nonce", nonce)

	http.Redirect(w, r, oauth2Config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
}

// LoggerMiddleware logs the incoming HTTP request & its duration.
func OIDCMiddleware(next http.Handler, config *models.Config) http.Handler {
	log.Println("INFO: registering OIDCMiddleware")

	context := context.Background()

	// Create OIDC provider
	provider, err := oidc.NewProvider(context, config.OIDC.ProviderUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create OAuth2 config
	oauth2Config := &oauth2.Config{
		ClientID:     config.OIDC.ClientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  strings.TrimSuffix(config.OIDC.RedirectApplication, "/") + callbackUrl,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Create OIDC config
	oidcConfig := &oidc.Config{
		ClientID: config.OIDC.ClientID,
	}

	// Create OIDC verifier
	verifier := provider.Verifier(oidcConfig)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// Fetch session cookie
			session, err := r.Cookie("session")
			if err != nil {
				// Redirect to login and return
				log.Printf("INFO: session cookie not found. Redirecting to login page. %s", err)
				loginRedirect(w, r, oauth2Config)
				return
			}

			// Validate session cookie
			idToken, err := verifier.Verify(context, session.Value)
			if err != nil {
				// Redirect to login and return
				log.Printf("INFO: session token has expired or is invalid. Redirecting to login page. %s", err)
				loginRedirect(w, r, oauth2Config)
				return
			}

			// Fetch claims from validated JWT
			IDTokenClaims := struct {
				Email string `json:"email"`
			}{}
			if err := idToken.Claims(&IDTokenClaims); err != nil {
				log.Printf("ERROR: could not fetch claims from vaidated jwt. %s", err)
				utils.ErrorJsonResponse(w, "something went wrong", http.StatusInternalServerError)
				return
			}

			// Set header with authenticated user email
			w.Header().Add("X-Auth-Email", IDTokenClaims.Email)

			next.ServeHTTP(w, r)
		},
	)
}
