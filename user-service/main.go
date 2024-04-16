package main

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/user-service/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Msg("Unable to load the env file")
	}

	// Replace this with your MongoDB Atlas connection string
	connectionString := os.Getenv("MONGO_URL")

	// Set MongoDB connection options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal().Msg("Internal server error: Unable to connect to Mongo")
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal().Msg("Internal server error: Unable to talk to Mongo")
	}

	logger.Info().Msg("Connected to Mongo!")

	// You can now use the "client" variable to interact with your MongoDB database.
	users := client.Database(os.Getenv("MONGO_DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION_USERS"))

	// Setup the router
	r := routes.SetupRouter(users, logger)
	logger.Info().Msg("Setup Complete. Starting user-service...")
	r.Run(":" + os.Getenv("PORT"))
}
