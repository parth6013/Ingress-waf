package middlewares

import (
	"armor/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// ElasticsearchClient initializes the connection.
var ElasticsearchClient *elasticsearch.Client

func InitElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"}, // Elasticsearch URL
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
	ElasticsearchClient = client
	log.Println("✅ Elasticsearch client initialized")
}

// LogEntry represents a log structure for Elasticsearch.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Status    int    `json:"status"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Duration  string `json:"duration"`
	Error     string `json:"error,omitempty"`
	Trace     string `json:"trace,omitempty"`
}

type LogEntryCoraza struct {
	Timestamp string `json:"timestamp"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
	Data      string `json:"data"`
	URI       string `json:"uri"`
}

// sendToElasticsearch sends logs to Elasticsearch.
func sendToElasticsearch(logEntry LogEntry) {
	data, err := json.Marshal(logEntry)
	if err != nil {
		log.Println("❌ Failed to serialize log entry:", err)
		return
	}

	res, err := ElasticsearchClient.Index("logs", bytes.NewReader(data))
	if err != nil {
		log.Println("❌ Failed to send log to Elasticsearch:", err)
		return
	}
	defer res.Body.Close()
	log.Println("✅ Log sent to Elasticsearch:", logEntry)
}

func SendToElasticSearchCoraza(logEntry LogEntryCoraza) {
	data, err := json.Marshal(logEntry)
	if err != nil {
		log.Println("❌ Failed to serialize log entry Coraza:", err)
		return
	}

	res, err := ElasticsearchClient.Index("coraza", bytes.NewReader(data))
	if err != nil {
		log.Println("❌ Failed to send log to Elasticsearch Coraza:", err)
		return
	}
	defer res.Body.Close()
	log.Println("✅ Log sent to Elasticsearch Coraza:", logEntry)
}

// LoggerMiddleware logs the request & sends logs to Elasticsearch only for 400 or 500 errors.
func LoggerMiddleware(next http.Handler) http.Handler {
	log.Println("INFO: registering LoggerMiddleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := wrapResponseWriter(w)

		defer func() {
			if err := recover(); err != nil {
				// Log 500 errors
				logEntry := LogEntry{
					Timestamp: time.Now().Format(time.RFC3339),
					Status:    http.StatusInternalServerError,
					Method:    r.Method,
					Path:      r.URL.Path,
					Duration:  time.Since(start).String(),
					Error:     "Internal Server Error",
					Trace:     string(debug.Stack()),
				}
				go sendToElasticsearch(logEntry) // Send asynchronously
				utils.ErrorJsonResponse(w, "something went wrong", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(wrapped, r)

		// Only send logs for 400 or 500 errors
		if wrapped.status >= 400 {
			logEntry := LogEntry{
				Timestamp: time.Now().Format(time.RFC3339),
				Status:    wrapped.status,
				Method:    r.Method,
				Path:      r.URL.Path,
				Duration:  time.Since(start).String(),
			}
			go sendToElasticsearch(logEntry) // Send asynchronously
		}

		log.Printf("INFO: status=%d method=%s path=%s duration=%s",
			wrapped.status, r.Method, r.URL.Path, time.Since(start))
	})
}

// // LoggerMiddleware logs the incoming HTTP request & its duration.
// func LoggerMiddleware(next http.Handler) http.Handler {
// 	log.Println("INFO: registering LoggerMiddleware")
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			start := time.Now()

// 			defer func() {
// 				if err := recover(); err != nil {
// 					utils.ErrorJsonResponse(w, "something went wrong", http.StatusInternalServerError)
// 					log.Printf(
// 						"FATAL: status=500 method=%s path=%s duration=%s error=%s. Trace=%s",
// 						r.Method, r.URL.EscapedPath(), time.Since(start), err, debug.Stack(),
// 					)
// 				}
// 			}()

// 			wrapped := wrapResponseWriter(w)

// 			next.ServeHTTP(wrapped, r)

// 			log.Printf(
// 				"INFO: status=%d method=%s path=%s duration=%s", wrapped.status,
// 				r.Method, r.URL.EscapedPath(), time.Since(start),
// 			)
// 		},
// 	)
// }
