package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/admin-service/models"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddCabs(c *gin.Context, cabs *mongo.Collection, logger zerolog.Logger) {
	// Initialise cab struct with default values
	cab := models.Cab{
		Location: models.GeoJSONPoint{
			Name:        "MoveInSync HQ",
			Coordinates: []float64{77.64344998289836, 12.912447107980537},
		},
		Status:      "Available",
		LastUpdated: time.Now(),
	}

	// Get JSON body from request
	if err := c.ShouldBindJSON(&cab); err != nil {
		logger.Error().Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Insert cab into database
	_, err := cabs.InsertOne(context.Background(), cab)
	if err != nil {
		logger.Error().Msg("Failed to insert cab into database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert cab into database"})
		return
	}

	// Return success
	logger.Info().Msg("Cab successfully added")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
