// event.go
package handlers

import (
	"encoding/json"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
	"log"
	"strconv"

	"MindMesh-Service/app/internal/models"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

var events []models.Event

func RegisterEventsRoutes(r *chi.Mux, client *mongo.Client) {
	r.Get("/api/events", getAllEventsHandler(client))
	r.Get("/api/events/{eventID}", getEventHandler(client))
	r.Put("/api/events/{eventID}", updateEventHandler(client))
	r.Delete("/api/events/{eventID}", deleteEventHandler(client))
	// r.Post("/api/events/{eventID}", createEventHandler(client))
}

func getAllEventsHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Return all events as JSON
		coll := client.Database("mindmesh").Collection("events")
		filter := bson.D{}
		// Retrieves documents that match the query filer
		var results []bson.M
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			fmt.Println(err)
		}
		if err := cursor.All(context.TODO(), &results); err != nil {
			log.Fatal(err)
		}
		fmt.Println(results)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func getEventHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "eventID")
		for _, event := range events {
			if strconv.Itoa(event.ID) == eventID {
				// Return the event as JSON
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(event)
				return
			}
		}
		// If event is not found, return a 404
		http.NotFound(w, r)
	}
}

func updateEventHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "eventID")
		var updatedEvent models.Event
		err := json.NewDecoder(r.Body).Decode(&updatedEvent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for i, event := range events {
			if strconv.Itoa(event.ID) == eventID {
				// Update the event
				events[i] = updatedEvent
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		// If event is not found, return a 404
		http.NotFound(w, r)
	}
}

func deleteEventHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "eventID")
		coll := client.Database("mindmesh").Collection("events")
		i, err := strconv.Atoi(eventID)
		if err != nil {
			// ... handle error
			panic(err)
		}
		filter := bson.D{{"id", i}}

		// Deletes the document with the specified ID
		result, err := coll.DeleteOne(context.TODO(), filter)
		if err != nil {
			fmt.Println("Error deleting event:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			// If event is not found, return a 404
			http.NotFound(w, r)
			return
		}

		// Event deleted successfully
		fmt.Println("Event deleted successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Event deleted successfully"))
	}
}