// notes.go
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

var notes[]models.Note

func RegisterNotesRoutes(r *chi.Mux, client *mongo.Client) {
	r.Get("/api/notes", getAllNotesHandler(client))
	r.Get("/api/notes/{noteID}", getNoteHandler(client))
	r.Put("/api/notes/{noteID}", updateNoteHandler(client))
	r.Delete("/api/notes/{noteID}", deleteNoteHandler(client))
	r.Post("/api/notes", addNoteHandler(client))
}

func getAllNotesHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Return all notes as JSON
		coll := client.Database("mindmesh").Collection("notes")
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

func getNoteHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteID := chi.URLParam(r, "noteID")
		for _, note := range notes {
			if strconv.Itoa(note.ID) == noteID {
				// Return the note as JSON
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(note)
				return
			}
		}
		// If note is not found, return a 404
		http.NotFound(w, r)
	}
}

func updateNoteHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteID := chi.URLParam(r, "noteID")
		var updatedNote models.Note
		err := json.NewDecoder(r.Body).Decode(&updatedNote)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for i, note := range notes {
			if strconv.Itoa(note.ID) == noteID {
				// Update the note
				notes[i] = updatedNote
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		// If note is not found, return a 404
		http.NotFound(w, r)
	}
}

func deleteNoteHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteID := chi.URLParam(r, "noteID")
		coll := client.Database("mindmesh").Collection("notes")
		i, err := strconv.Atoi(noteID)
		if err != nil {
			// ... handle error
			panic(err)
		}
		filter := bson.D{{"id", i}}

		// Deletes the document with the specified ID
		result, err := coll.DeleteOne(context.TODO(), filter)
		if err != nil {
			fmt.Println("Error deleting note:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			// If note is not found, return a 404
			http.NotFound(w, r)
			return
		}

		// Note deleted successfully
		fmt.Println("Note deleted successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Note deleted successfully"))
	}
}

func addNoteHandler(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		noteID := chi.URLParam(r, "noteID")
		coll := client.Database("mindmesh").Collection("notes")
		i, err := strconv.Atoi(noteID)
		if err != nil {
			// ... handle error
			panic(err)
		}
		filter := bson.D{{"id", i}}

		// Deletes the document with the specified ID
		result, err := coll.DeleteOne(context.TODO(), filter)
		if err != nil {
			fmt.Println("Error deleting note:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if result.DeletedCount == 0 {
			// If note is not found, return a 404
			http.NotFound(w, r)
			return
		}

		// Note deleted successfully
		fmt.Println("Note deleted successfully")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Note deleted successfully"))
	}
}