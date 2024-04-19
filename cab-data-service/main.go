package main

import (
	"context"
	"os"
	"time"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Cab represents the data structure for a cab
type Cab struct {
	ID          string       `json:"id" bson:"_id"`
	Location    GeoJSONPoint `json:"location" bson:"location"`
	Status      string       `json:"status" bson:"status"`
	LastUpdated time.Time    `json:"last_updated" bson:"last_updated"`
}

// GeoJSONPoint represents a GeoJSON Point
type GeoJSONPoint struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates"`
}

var cabClients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Cab)

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
	cabs := client.Database(os.Getenv("MONGO_DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION_CABS"))

	// Setup the router
	router := gin.Default()

	// Setup Prometheus metrics
	p := ginprom.New(
		ginprom.Engine(router),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	router.Use(p.Instrument())

	// WebSocket handler for cabs
	router.POST("/cab/ws", func(c *gin.Context) {
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to upgrade WebSocket connection")
			return
		}
		defer conn.Close()

		// Add cab client to the cabClients map
		cabClients[conn] = true
		logger.Info().Msg("New WebSocket connection established for Cab data...")

		// Read from WebSocket connection
		for {
			var cab Cab
			if err := conn.ReadJSON(&cab); err != nil {
				logger.Error().Msg("Failed to read message from WebSocket connection")
				delete(cabClients, conn)
				break
			}
			// Broadcast cab location and status updates to all cab clients
			broadcast <- cab
		}
	})

	// Start WebSocket server for cabs
	go startCabServer(cabs, logger)
	logger.Info().Msg("WebSocket server started for Cab data...")
	router.Run(":" + os.Getenv("PORT"))
}

func startCabServer(cabs *mongo.Collection, logger zerolog.Logger) {
	for {
		// Wait for cab location and status updates
		cab := <-broadcast

		// Update or insert cab location and status data into the database
		filter := bson.M{"_id": cab.ID}
		update := bson.M{
			"$set": bson.M{
				"location":     cab.Location,
				"status":       cab.Status,
				"last_updated": time.Now(),
			},
		}
		opts := options.Update().SetUpsert(true)
		_, err := cabs.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			logger.Error().Msg("Internal server error: Unable to update cab data")
			continue
		}
	}
}
