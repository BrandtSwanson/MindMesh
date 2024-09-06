// personal.go
package handlers

import (
	"encoding/json"
	"net/http"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"fmt"
	"log"
	"time"

	"MindMesh-Service/app/internal/models"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
)

var personal []models.Personal

func RegisterPersonalRoutes(r *chi.Mux, client *mongo.Client) {
	r.Post("/api/personal", createPersonalHandler(client))
	r.Get("/api/personal", getAllPersonalHandler(client))
	// Add other note-related routes
}

func createPersonalHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newPersonalData models.Personal
		err := json.NewDecoder(r.Body).Decode(&newPersonalData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newPersonalData.TimeStamp = time.Now()

		// Assign a unique ID and add the note to the slice
		personal = append(personal, newPersonalData)

		fmt.Println("HERE")
		fmt.Printf("%#v", newPersonalData)

		// Return the created note as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newPersonalData)

		coll := client.Database("mindmesh").Collection("personal_data_timeseries")
		result, err := coll.InsertOne(context.TODO(), newPersonalData)
		fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	}
}

func getAllPersonalHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Return all notes as JSON
		coll := client.Database("mindmesh").Collection("personal_data_timeseries")
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