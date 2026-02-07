package main

import (
	"fmt"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/ASR11104/the-mallu-hangman/internal/config"
	"github.com/ASR11104/the-mallu-hangman/internal/handlers"
)

func main() {
	godotenv.Load() // only in dev

	mux := http.NewServeMux()
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/health", handlers.Health)
	mux.Handle("/movie", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.Movies(w, r, cfg)
	}))
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Server is running on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
