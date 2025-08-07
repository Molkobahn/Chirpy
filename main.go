package main

import(
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"database/sql"
	"github.com/molkobahn/Chirpy/internal/database"
	"encoding/json"
	"github.com/joho/godotenv"
)


type apiConfig struct {
	fileserverHits	atomic.Int32
	db	*database.Queries
	platform string
}


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}


func (cfg *apiConfig) hitsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	hits := cfg.fileserverHits.Load()
	hitsString := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)
	w.Write([]byte(hitsString))
}


func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden", nil)
		return
	}
	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete users", err)
	}
	cfg.fileserverHits.Swap(0)
}

func (cfg *apiConfig)createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	log.Printf("Struct of Parameters: %v", params)
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
		} 
	newUser := mapUser(user)
	respondWithJSON(w, http.StatusCreated, newUser)
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
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHitsHandler)
	
	server := &http.Server {
	Addr:		":" + port,
	Handler:	mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
