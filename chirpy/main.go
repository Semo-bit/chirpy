package main
import (
    "log"
    "net/http"
)

func main() {
    const port = "8080"
    mux := http.NewServeMux()
    srv := &http.Server{
        Addr: ":" + port,
        Handler: mux,
    }
    mux.HandleFunc("/healthz", healtzHandler)

    fileServer := http.FileServer(http.Dir("."))
    mux.Handle("/app/", http.StripPrefix("/app", fileServer))
    log.Printf("Serving on port: %s\n", port)
    log.Fatal(srv.ListenAndServe())
}

func healtzHandler(w http.ResponseWriter, r *http.Request) {
    // Set the Content-Type header
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    
    // Write the status code
    w.WriteHeader(http.StatusOK)
    
    // Write the body text
    w.Write([]byte("OK"))
    
}
