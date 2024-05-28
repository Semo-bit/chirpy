package main

import (
	"fmt"
	"log"
	"net/http"
)

// Struct to hold the number of times our page was visted
type apiConfig struct {
    fileServerHits int
}

// Hanlder for health check
func healtzHandler(w http.ResponseWriter, r *http.Request) {
    // Set the Content-Type header
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    
    // Write the status code
    w.WriteHeader(http.StatusOK)
    
    // Write the body text
    w.Write([]byte("OK"))
    
}

// Function to count the number of hits
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileServerHits++
        next.ServeHTTP(w, r)
    })
} 

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hits: %d", cfg.fileServerHits)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
    cfg.fileServerHits = 0
    w.WriteHeader(http.StatusOK)
}

func main() {
    const port = "8080"
    mux := http.NewServeMux()
    srv := &http.Server{
        Addr: ":" + port,
        Handler: mux,
    }
    mux.HandleFunc("GET /healthz", healtzHandler)

    // initialzing the apiConfig
    apiCfg := &apiConfig{fileServerHits: 0}
    
    fileServer := http.FileServer(http.Dir("."))
    handler := http.StripPrefix("/app", fileServer)
    mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

    // creating a new handler to report hits
    mux.HandleFunc("GET /metrics", apiCfg.metricsHandler)
    // creating a new handler to reset hits
    mux.HandleFunc("/reset", apiCfg.resetHandler)

    log.Printf("Serving on port: %s\n", port)
    log.Fatal(srv.ListenAndServe())
}  
