package main

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/admin-service/routes"
	"go.mongodb.org/mongo-driver/bson"
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
	admins := client.Database(os.Getenv("MONGO_DB1_NAME")).Collection(os.Getenv("MONGO_COLLECTION_ADMINS"))
	cabs := client.Database(os.Getenv("MONGO_DB2_NAME")).Collection(os.Getenv("MONGO_COLLECTION_CABS"))

	// Create a 2dsphere index on the "location" field of the "cabs" collection
	_, err = cabs.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.M{
				"location": "2dsphere",
			},
		},
	)
	if err != nil {
		logger.Fatal().Msg("Internal server error: Unable to create 2dsphere index on 'location' field of 'cabs' collection")
	}
	logger.Info().Msg("2dsphere index created on 'location' field of 'cabs' collection")

	// Setup the router
	r := routes.SetupRouter(admins, cabs, logger)
	logger.Info().Msg("Setup Complete. Starting user-service...")
	r.Run(":" + os.Getenv("PORT"))
}
