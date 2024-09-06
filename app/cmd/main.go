package main

import (
	"log"
	"net/http"
	"context"

	"MindMesh-Service/app/internal/database"
	"MindMesh-Service/app/internal/handlers"
	"MindMesh-Service/app/internal/middleware"

	"github.com/go-chi/chi"
)

func main() {
	// MongoDB connection
	client := database.ConnectMongoDB()
	defer client.Disconnect(context.TODO())

	// Create a new Chi router
	r := chi.NewRouter()

	// Apply CORS middleware
	r.Use(middleware.NewCORS().Handler)

	// Notes endpoints
	handlers.RegisterNotesRoutes(r, client)

	// Events endpoints
	handlers.RegisterEventsRoutes(r, client)

	// Personal data endpoints
	handlers.RegisterPersonalRoutes(r, client)

	// Start server
	log.Fatal(http.ListenAndServe(":8181", r))
}
