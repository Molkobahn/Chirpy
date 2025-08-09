package main

import(
	"net/http"
	"log"
	"sync/atomic"
	_ "github.com/lib/pq"
	"os"
	"database/sql"
	"github.com/molkobahn/Chirpy/internal/database"
	"github.com/joho/godotenv"
)


type apiConfig struct {
	fileserverHits	atomic.Int32
	db	*database.Queries
	platform string
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
    	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}



func main() {
	const filepathRoot = "."
	const port = "8080"
	godotenv.Load()
	// Setup database connection
	dbURL := os.Getenv("DB_URL")
	log.Printf("The URL: %v", dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open Database: %v", err)
	}
	dbQueries := database.New(db)
	// Setup config struct
	apiCfg := new(apiConfig)
	apiCfg.db = dbQueries
	platform := os.Getenv("PLATFORM")
	apiCfg.platform = platform 
	mux := http.NewServeMux()
	
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.chirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.loginHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHitsHandler)
	
	server := &http.Server {
	Addr:		":" + port,
	Handler:	mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
