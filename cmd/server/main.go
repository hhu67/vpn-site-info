package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"vpn-site-info/internal/config"
	"vpn-site-info/internal/database"
	"vpn-site-info/internal/handlers"
	"vpn-site-info/internal/middleware"

	"github.com/joho/godotenv"
)

func init() {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[%s] START %s %s from %s", r.Method, r.RequestURI, time.Now().Format("2006-01-02 15:04:05"), r.RemoteAddr)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Printf("[%s] END %s completed in %v", r.Method, r.RequestURI, duration)
	})
}

func serveSPA(distPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[SPA] Serving request for path: %s", r.URL.Path)
		filePath := filepath.Join(distPath, r.URL.Path)

		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			log.Printf("[SPA] File found: %s, serving", filePath)
			http.FileServer(http.Dir(distPath)).ServeHTTP(w, r)
			return
		}

		log.Printf("[SPA] File not found: %s, serving index.html", filePath)
		http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
	}
}

func main() {
	log.Println("=== VPN Site Info Server Starting ===")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()
	log.Printf("Config loaded - Port: %s, JWT Secret: %d chars", cfg.Port, len(cfg.JWTSecret))

	log.Println("Connecting to database...")
	db, err := database.Connect(cfg.PSQL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully")
	defer db.Close()

	log.Println("Initializing database schema...")
	if err := database.InitSchema(db); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}
	log.Println("Database schema initialized")

	h := handlers.New(db, cfg.JWTSecret)
	log.Println("Handlers initialized")

	mux := http.NewServeMux()

	log.Println("Registering API routes...")
	mux.HandleFunc("/api/login", h.Login)
	mux.HandleFunc("/api/check-password-exists", h.CheckPasswordExists)
	mux.HandleFunc("/api/create-password", h.CreatePassword)

	mux.Handle("/api/insert/vpn", middleware.AuthMiddleware(http.HandlerFunc(h.InsertVPN), cfg.JWTSecret))
	mux.Handle("/api/update/vpn", middleware.AuthMiddleware(http.HandlerFunc(h.UpdateVPN), cfg.JWTSecret))
	mux.Handle("/api/delete/vpn", middleware.AuthMiddleware(http.HandlerFunc(h.DeleteVPN), cfg.JWTSecret))
	mux.Handle("/api/list/vpn", middleware.AuthMiddleware(http.HandlerFunc(h.ListVPN), cfg.JWTSecret))
	mux.Handle("/api/change-password", middleware.AuthMiddleware(http.HandlerFunc(h.ChangePassword), cfg.JWTSecret))
	log.Println("API routes registered")

	log.Println("Setting up SPA routes...")
	mux.HandleFunc("/", serveSPA("./frontend/dist"))

	handler := loggingMiddleware(mux)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("=== Server starting on port %s ===", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
