package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/user-service/models"
	"github.com/sankalp-12/moveinsync/user-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// BookTrip handles the API endpoint for booking a trip
func BookTrip(c *gin.Context, cabs *mongo.Collection, logger zerolog.Logger) {
	// Get JSON body from request
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		logger.Error().Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Prepare the URL with query parameters
	url := "http://admin:8081/api/v1/cab/available?latitude=" + location.Latitude + "&longitude=" + location.Longitude

	resp, err := utils.SendGetRequest(url, logger)
	if err != nil {
		logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		logger.Error().Msg("Failed to fetch available cabs")
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch available cabs"})
		return
	}

	// Decode the response body
	var results []bson.M
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		logger.Error().Msg("Failed to decode response body")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response body"})
		return
	}

	// Return the results
	logger.Info().Msg("Successfully fetched best available cab")
	c.JSON(http.StatusOK, results)
}

// DisplayNearbyCabs handles the API endpoint for displaying nearby busy cabs
func DisplayNearbyCabs(c *gin.Context, cabs *mongo.Collection, logger zerolog.Logger) {
	// Get JSON body from request
	var location models.Location
	if err := c.ShouldBindJSON(&location); err != nil {
		logger.Error().Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Prepare the URL with query parameters
	url := "http://admin:8081/api/v1/cab/busy?latitude=" + location.Latitude + "&longitude=" + location.Longitude

	resp, err := utils.SendGetRequest(url, logger)
	if err != nil {
		logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		logger.Error().Msg("Failed to fetch nearby busy cabs")
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch nearby busy cabs"})
		return
	}

	// Decode the response body
	var results []bson.M
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		logger.Error().Msg("Failed to decode response body")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response body"})
		return
	}

	// Return the results
	logger.Info().Msg("Successfully fetched nearby busy cabs")
	c.JSON(http.StatusOK, results)
}
