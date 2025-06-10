package main

import (
	"encoding/json"
	"expvar"
	"flag"

	_ "expvar"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mikebellcoder/nlp"
	"github.com/mikebellcoder/nlp/stemmer"
)

var (
	stemCalls = expvar.NewInt("stem.calls")
)

var config struct {
	Addr string
}

func main() {
	time.Local = time.UTC
	// config
	config.Addr = os.Getenv("NLP_ADDR")
	if config.Addr == "" {
		config.Addr = ":8080"
	}
	flag.StringVar(&config.Addr, "addr", config.Addr, "HTTP server address")
	flag.Parse()

	// set time to UTC
	// os.Setenv("TZ", "UTC")
	// init
	api := &API{
		log: slog.Default().With("app", "nlp"),
	}
	// routing
	http.HandleFunc("GET /health", api.middlewareLogger(api.healthHandler))
	http.HandleFunc("POST /tokenize", api.middlewareLogger(api.tokenizeHandler))
	http.HandleFunc("GET /stem/{word}", api.middlewareLogger(api.stemHandler))

	api.log.Info("server starting", "addr", config.Addr)
	if err := http.ListenAndServe(config.Addr, nil); err != nil {
		fmt.Fprintf(os.Stderr, "failed to start server: %v\n", err)
		os.Exit(1)
	}
}

type API struct {
	log *slog.Logger
}

func (a *API) healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := health(); err != nil {
		a.log.Error("health", "error", err)
		http.Error(w, "health check failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

func health() error {
	// Simulate a health check
	return nil
}

// Middleware to log requests
func (a *API) middlewareLogger(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			msg := fmt.Sprintf("Received request: %s %s processed in %s", r.Method, r.URL.Path, duration)
			a.log.Info(msg)
		}()

		next(w, r)
	}
}

func (a *API) stemHandler(w http.ResponseWriter, r *http.Request) {
	stemCalls.Add(1)
	word := r.PathValue("word")
	a.log.Info("stem", "word", word)
	fmt.Fprintln(w, stemmer.Stem(word))
}

func (a *API) tokenizeHandler(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("read", "error", err, "remote", r.RemoteAddr)
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	text := string(data)
	if len(text) == 0 {
		a.log.Error("read", "error", "empty request")
		http.Error(w, "empty request", http.StatusBadRequest)
		return
	}

	tokens := nlp.Tokenize(text)

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"text":   text,
		"tokens": tokens,
	}
	json.NewEncoder(w).Encode(&resp)
}
