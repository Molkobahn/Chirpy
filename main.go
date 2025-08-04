package main

import(
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
)


type apiConfig struct {
	fileserverHits	atomic.Int32
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}


func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	hitsString := fmt.Sprintf("Hits: %v", hits)
	w.Write([]byte(hitsString))
}


func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Swap(0)
}


func healthzHandler(w http.ResponseWriter, r *http.Request) {
    	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}


func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := new(apiConfig)

	mux := http.NewServeMux()
	
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/reset", apiCfg.resetHitsHandler)
	
	server := &http.Server {
	Addr:		":" + port,
	Handler:	mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
