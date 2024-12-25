//go:generate go run github.com/99designs/gqlgen
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	// Load environment variables from the .env file explicitly
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables explicitly and set defaults if needed
	accountURL := os.Getenv("ACCOUNT_SERVICE_URL")
	catalogURL := os.Getenv("CATALOG_SERVICE_URL")
	orderURL := os.Getenv("ORDER_SERVICE_URL")

	// If the environment variables are not found, set default values
	if accountURL == "" {
		accountURL = "http://localhost:8081"
	}
	if catalogURL == "" {
		catalogURL = "http://localhost:8082"
	}
	if orderURL == "" {
		orderURL = "http://localhost:8083"
	}

	// Print out the URLs for debugging
	log.Printf("Account URL: %s, Catalog URL: %s, Order URL: %s", accountURL, catalogURL, orderURL)

	// Initialize AppConfig with values
	var cfg AppConfig
	cfg.AccountURL = accountURL
	cfg.CatalogURL = catalogURL
	cfg.OrderURL = orderURL

	// Check if any of the required URLs are still empty
	if cfg.AccountURL == "" || cfg.CatalogURL == "" || cfg.OrderURL == "" {
		log.Fatal("One or more environment variables are missing or empty")
	}

	// Initialize GraphQL server (assuming NewGraphQLServer and ToExecutableSchema are defined)
	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/graphql", handler.GraphQL(s.ToExecutableSchema()))
	http.Handle("/playground", handler.Playground("GraphQL Playground", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
