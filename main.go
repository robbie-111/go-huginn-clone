package main

import (
	"log"
	"net/http"
	"os"

	"go-huginn-clone/router"
)

func main() {
	r := router.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	addr := ":" + port
	log.Printf("Huginn (Go) starting on http://localhost%s", addr)
	log.Printf("Login with any username/password combination")
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
