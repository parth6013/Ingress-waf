package handlers

import (
	"armor/models"
	"armor/utils"
	"context"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

const (
	CallbackHandlerUrl = "/auth/oidc/callback"
)

var (
	// clientSecret = os.Getenv("DUO_OAUTH2_CLIENT_SECRET")
	clientSecret = "PU9ug0fnSl3fTUdPWFlPWEq7sXigTgqp"
	test_url     = "http://localhost:3000/auth/oidc/callback"
)

type OIDCHandler struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	ctx          context.Context
}

func NewOIDCHandler(config *models.Config) *OIDCHandler {
	ctx := context.Background()

	// Create OIDC provider
	provider, err := oidc.NewProvider(ctx, config.OIDC.ProviderUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create OAuth2 config
	oauth2Config := &oauth2.Config{
		ClientID:     config.OIDC.ClientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		// RedirectURL:  strings.TrimSuffix(config.OIDC.RedirectApplication, "/") + strings.TrimPrefix(CallbackHandlerUrl, "GET "),
		RedirectURL: test_url,
		Scopes:      []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Create OIDC config
	oidcConfig := &oidc.Config{
		ClientID: config.OIDC.ClientID,
	}

	// Create OIDC verifier
	verifier := provider.Verifier(oidcConfig)

	return &OIDCHandler{
		oauth2Config: oauth2Config,
		verifier:     verifier,
		ctx:          ctx,
	}
}

func (o *OIDCHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Compare the value of state parameter in the request with the value of state cookie
	state, err := r.Cookie("state")
	if err != nil {
		log.Printf("WARN: state cookie not found")
		utils.ErrorJsonResponse(w, "state not found", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		log.Printf("WARN: the value of state from the cookie does not match the value of state from the callback request")
		utils.ErrorJsonResponse(w, "invalid state", http.StatusBadRequest)
		return
	}

	// exchange the code sent back to the callback URL for an OAuth token
	oauth2Token, err := o.oauth2Config.Exchange(o.ctx, r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("ERROR: failed to exchange id_token with the auth code. %s", err)
		utils.ErrorJsonResponse(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	// from oauth2Token, fetch the rawIDToken ( which is the JWT token ) from the field id_token.
	// This token can be saved as a session token on the client side
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Printf("ERROR: no id_token field in oauth2 token. %s", err)
		utils.ErrorJsonResponse(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	// verify id_token sent to the callback URL
	idToken, err := o.verifier.Verify(o.ctx, rawIDToken)
	if err != nil {
		log.Printf("WARN: failed to verify id_token. %s", err)
		utils.ErrorJsonResponse(w, "failed to verify ID Token", http.StatusUnauthorized)
		return
	}

	// Compare the value of nonce parameter in the request with the value of nonce cookie
	nonce, err := r.Cookie("nonce")
	if err != nil {
		log.Printf("WARN: nonce cookie not found")
		utils.ErrorJsonResponse(w, "nonce not found", http.StatusBadRequest)
		return
	}
	if idToken.Nonce != nonce.Value {
		log.Printf("WARN: the value of nonce from the cookie does not match the value of state from the id_token")
		utils.ErrorJsonResponse(w, "invalid nonce", http.StatusBadRequest)
		return
	}

	// Set session cookie
	utils.SetCookie(w, r, "session", rawIDToken)

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}
