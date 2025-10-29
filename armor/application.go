package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"armor/handlers"
	"armor/middlewares"
	"armor/models"
	"armor/utils"
)

type Application struct {
	config   *models.Config
	handlers http.Handler
}

func NewApplication() *Application {
	config, err := utils.LoadConfigs()
	if err != nil {
		log.Fatalf("FATAL: Failed to load config. %s", err)
	}

	application := &Application{config: config}
	application.registerHandlers()
	application.registerMiddlewares()

	return application
}

// Register Handlers
func (a *Application) registerHandlers() {
	router := http.NewServeMux()
	router.HandleFunc(handlers.HealthHandlerUrl, handlers.HealthHandler)

	// Proxy handler
	// Handle the URL validation in the config parser
	proxyHandler := handlers.NewProxyHandler(a.config)

	if a.config.OIDC.Enable {
		// Register OIDC callback auth url
		oidcHandler := handlers.NewOIDCHandler(a.config)
		router.HandleFunc(handlers.CallbackHandlerUrl, oidcHandler.CallbackHandler)

		// Create a protected router for all proxy requests
		protectedRouter := http.NewServeMux()
		protectedRouter.Handle("/", proxyHandler)

		// Register OIDC middleware to the protected router
		oidcMiddleware := middlewares.OIDCMiddleware(protectedRouter, a.config)

		// Register OIDC middleware to the main router
		router.Handle("/", oidcMiddleware)
	} else {
		// Register proxy handler to the main router
		router.Handle("/", proxyHandler)
	}

	a.handlers = router

	log.Println("INFO: completed registering handlers")
}

// Register Middlewares
func (a *Application) registerMiddlewares() {
	// Initialize Elasticsearch if enabled
	if a.config.Elasticsearch.Enable {
		log.Println("INFO: Elasticsearch enabled, initializing client")
		middlewares.InitElasticsearch(a.config.Elasticsearch.URL)
	} else {
		log.Println("INFO: Elasticsearch disabled, logs will be printed to console only")
	}

	// Register logger middleware
	a.handlers = middlewares.LoggerMiddleware(a.handlers)

	// Register coraza middleware
	a.handlers = middlewares.CorazaMiddleware(a.handlers, a.config)

	log.Println("INFO: completed registering middlewares")
}

// Serve starts the http proxy server
func (a *Application) Serve() {
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", a.config.Server.Host, a.config.Server.Port),
		Handler:      a.handlers,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownChan := make(chan error)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigChan

		log.Printf("WARN: shutting down server. Signal: %s", sig.String())
		ctx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownRelease()

		shutdownChan <- server.Shutdown(ctx)
	}()

	log.Printf("INFO: starting Armor server on: %s", server.Addr)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("ERROR: server failed to start. %s", err)
	}

	err = <-shutdownChan
	if err != nil {
		log.Fatalf("ERROR: error while shutting down server. %s", err)
	}

	log.Println("INFO: server shutdown gracefully")
}
