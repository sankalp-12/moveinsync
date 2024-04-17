package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SuggestAvailableCabs handles the API endpoint for allocating available cabs
func SuggestAvailableCabs(c *gin.Context, cabs *mongo.Collection, logger zerolog.Logger) {
	// Parse latitude and longitude from query parameters
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")

	// Convert latitude and longitude to float64
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		logger.Error().Err(err).Msg("Invalid latitude")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}
	lng, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		logger.Error().Err(err).Msg("Invalid longitude")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	// Specify the coordinates of the point
	point := bson.A{lng, lat}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up the $geoNear stage options
	geoNearStage := bson.D{
		{Key: "$geoNear", Value: bson.D{
			{Key: "near", Value: bson.D{
				{Key: "type", Value: "Point"},
				{Key: "coordinates", Value: point},
			}},
			{Key: "distanceField", Value: "distance"},
			{Key: "spherical", Value: true},
		}},
	}

	// Define the $match stage to filter cabs with status "Available"
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "status", Value: "Available"},
		}},
	}

	// Define the $limit stage to limit the number of results to 1
	limitStage := bson.D{
		{Key: "$limit", Value: 1},
	}

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{geoNearStage, matchStage, limitStage}

	// Execute the aggregation pipeline
	cursor, err := cabs.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to execute aggregation pipeline")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute aggregation pipeline"})
		return
	}
	defer cursor.Close(ctx)

	// Iterate through the results and handle them as needed
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		logger.Error().Err(err).Msg("Failed to iterate through aggregation results")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate through aggregation results"})
		return
	}

	// Return the results
	c.JSON(http.StatusOK, results)
}

// SuggestBusyCabs handles the API endpoint for suggesting nearby busy cabs
func SuggestBusyCabs(c *gin.Context, cabs *mongo.Collection, logger zerolog.Logger) {
	// Parse latitude and longitude from query parameters
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")

	// Convert latitude and longitude to float64
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		logger.Error().Err(err).Msg("Invalid latitude")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}
	lng, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		logger.Error().Err(err).Msg("Invalid longitude")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	// Specify the coordinates of the point
	point := bson.A{lng, lat}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up the $geoNear stage options
	geoNearStage := bson.D{
		{Key: "$geoNear", Value: bson.D{
			{Key: "near", Value: bson.D{
				{Key: "type", Value: "Point"},
				{Key: "coordinates", Value: point},
			}},
			{Key: "distanceField", Value: "distance"},
			{Key: "spherical", Value: true},
			{Key: "maxDistance", Value: 5000},
		}},
	}

	// Define the $match stage to filter cabs with status "Available"
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "status", Value: "Busy"},
		}},
	}

	// Define the $limit stage to limit the number of results to 1
	limitStage := bson.D{
		{Key: "$limit", Value: 5},
	}

	// Define the aggregation pipeline
	pipeline := mongo.Pipeline{geoNearStage, matchStage, limitStage}

	// Execute the aggregation pipeline
	cursor, err := cabs.Aggregate(ctx, pipeline)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to execute aggregation pipeline")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute aggregation pipeline"})
		return
	}
	defer cursor.Close(ctx)

	// Iterate through the results and handle them as needed
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		logger.Error().Err(err).Msg("Failed to iterate through aggregation results")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate through aggregation results"})
		return
	}

	// Return the results
	c.JSON(http.StatusOK, results)
}
