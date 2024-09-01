package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"MindMesh-Service/internal/middleware" // Replace with your actual module name

	"github.com/go-chi/chi"
)

// ExampleResponse is a struct representing the JSON response.
type ExampleResponse struct {
	Sentence string `json:"sentence"`
}

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Goal struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Event struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
}

type Personal struct {
	TimeStamp time.Time `json:"timestamp"`
	ID      int       `json:"id"`
	Water      int       `json:"water"`
	Joyousness int `json:"joyousness`
	Relaxation int `json:"relaxation`
	Alertness int `json:"alertness`
	ScreenTime int `json:"screenTime`
	Satisfaction int `json:"satisfaction`
	OutsideTime int `json:"outsideTime`
};

func main() {
	connectionString := "mongodb://192.168.4.39:27017/database"

	// // Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Create a MongoDB client
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context with a timeout (adjust the timeout as needed)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB server
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Connected to MongoDB!")
	
	//Check for database collections, instantiate ones that don't exist

	// Now you can perform database operations using the 'client'

	// // Don't forget to close the connection when you're done
	// err = client.Disconnect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Create a new Chi router
	r := chi.NewRouter()
	var notes []Note
	var events []Event
	var personal []Personal

	// Apply your custom CORS middleware to the router
	r.Use(middleware.NewCORS().Handler) // Use your custom CORS middleware

	// Retrieve all notes
	r.Get("/api/notes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ayo----\n\nBo\n-----")
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
	})

	// Retrieve a single note by ID
	r.Get("/api/notes/{noteID}", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Update a note by ID
	r.Put("/api/notes/{noteID}", func(w http.ResponseWriter, r *http.Request) {
		noteID := chi.URLParam(r, "noteID")
		var updatedNote Note
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
	})

	// Delete a note by ID
	r.Delete("/api/notes/{noteID}", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Create a new note
	r.Post("/api/notes", func(w http.ResponseWriter, r *http.Request) {

		var newNote Note
		err := json.NewDecoder(r.Body).Decode(&newNote)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Assign a unique ID and add the note to the slice
		notes = append(notes, newNote)

		// Return the created note as JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newNote)

		coll := client.Database("mindmesh").Collection("notes")
		result, err := coll.InsertOne(context.TODO(), newNote)
		fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	})

	// Create a new personal entry
	r.Post("/api/personal", func(w http.ResponseWriter, r *http.Request) {

		var newPersonalData Personal
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
	})

	// Retrieve all notes
	r.Get("/api/personal", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ayo----\n\nBo\n-----")
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
	})

	// Retrieve all events
	r.Get("/api/events", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Retrieve a single event by ID
	r.Get("/api/events/{eventID}", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Update a event by ID
	r.Put("/api/events/{eventID}", func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "eventID")
		var updatedEvent Event
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
	})

	// Delete a event by ID
	r.Delete("/api/events/{eventID}", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Start the server
	http.ListenAndServe(":8080", r)

}
