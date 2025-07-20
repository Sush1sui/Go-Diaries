package main

import (
	"context"
	"fmt"
	"gocrudapp/repository/mongodb"
	"gocrudapp/usecase"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Initialize the .env file and MongoDB connection
func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	slog.Info(".env loaded successfully")
}
func mongoConnection() *mongo.Client {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
  client, err := mongo.Connect(opts)
  if err != nil {
    panic(err)
  }

  // Send a ping to confirm a successful connection
  if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
    panic(err)
  }
  fmt.Println("DB Connected!")

	return client
}

func main() {
	mongoClient := mongoConnection() // establish MongoDB connection
	defer mongoClient.Disconnect(context.Background()) // ensure the connection is closed when done

	// mongo connection to user collection
	userCollection := mongoClient.Database(os.Getenv("MONGODB_NAME")).Collection(os.Getenv("MONGODB_COLLECTION"))

	// UserService instance
	userService := usecase.UserService{
		DBClient: mongodb.MongoClient{
			Client: *userCollection,
		},
	}

	r := chi.NewRouter() // create new instance of chi router
	r.Use(middleware.Logger) // use the Logger middleware to log requests
	r.Use(middleware.SetHeader("Content-Type", "application/json")) // set the Content-Type header for all responses

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to the Go CRUD app!"))
	})

	// define route
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Server is healthy!"))
		})
		r.Post("/users", userService.CreateUser)
		r.Get("/users/{id}", userService.GetUserByID)
		r.Get("/users", userService.GetAllUsers)
		r.Put("/users/{id}", userService.UpdateUserAgeByID)
		r.Delete("/users/{id}", userService.DeleteUserByID)
		r.Delete("/users", userService.DeleteAllUsers)
	})

	slog.Info("Starting server on port 4444")
	http.ListenAndServe(":4444", r)
}