package main

import (
	"log"
	"net/http"

	"vpn-site-info/internal/config"
	"vpn-site-info/internal/database"
	"vpn-site-info/internal/handlers"
	"vpn-site-info/internal/middleware"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	db, err := database.Connect(cfg.PSQL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	if err := database.InitSchema(db); err != nil {
		log.Fatal("Failed to initialize schema:", err)
	}

	h := handlers.New(db, cfg.JWTSecret)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", h.Login)
	mux.HandleFunc("/api/check-password-exists", h.CheckPasswordExists)
	mux.HandleFunc("/api/create-password", h.CreatePassword)

	mux.Handle("/api/insert/vpn", middleware.AuthMiddleware(http.HandlerFunc(h.InsertVPN), cfg.JWTSecret))
	mux.Handle("/api/update/vpn", middleware.AuthMiddleware(http.HandlerFunc(h.UpdateVPN), cfg.JWTSecret))
	mux.Handle("/api/delete/vpn", middleware.AuthMiddleware(http.HandlerFunc(h.DeleteVPN), cfg.JWTSecret))
	mux.Handle("/api/list/vpn", middleware.AuthMiddleware(http.HandlerFunc(h.ListVPN), cfg.JWTSecret))
	mux.Handle("/api/change-password", middleware.AuthMiddleware(http.HandlerFunc(h.ChangePassword), cfg.JWTSecret))

	fs := http.FileServer(http.Dir("./frontend/dist"))
	mux.Handle("/", fs)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed:", err)
	}
}
