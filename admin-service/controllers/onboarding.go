package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/admin-service/models"
	"github.com/sankalp-12/moveinsync/admin-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(c *gin.Context, admins *mongo.Collection, logger zerolog.Logger) {
	// Get JSON body from request
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		logger.Error().Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the username already exists
	var existingAdmin models.Admin
	err := admins.FindOne(context.TODO(), bson.M{"username": admin.Username}).Decode(&existingAdmin)
	if err == nil {
		// Username already exists
		logger.Error().Msg("Username already exists")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		// Error occurred while querying the database
		logger.Error().Msg("Error querying the database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
		return
	}

	// Hash the password
	hash, err := utils.HashPassword(admin.Password)
	if err != nil {
		logger.Error().Msg("Unable to hash password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to hash password"})
		return
	}

	adminModel := models.Admin{
		Username: admin.Username,
		Password: hash,
	}

	// Insert user in MongoDB
	_, err = admins.InsertOne(context.TODO(), adminModel)
	if err != nil {
		logger.Error().Msg("Unable to insert admin in database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to insert admin in database"})
		return
	}

	// Return success
	logger.Info().Msg("Admin successfully created")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func Login(c *gin.Context, admins *mongo.Collection, logger zerolog.Logger) {
	// Get JSON body from request
	var adminRequest models.Admin
	if err := c.ShouldBindJSON(&adminRequest); err != nil {
		logger.Error().Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find the user in MongoDB
	var adminDB models.Admin
	err := admins.FindOne(context.TODO(), bson.M{"username": adminRequest.Username}).Decode(&adminDB)
	if err != nil {
		logger.Error().Msg("Username is incorrect")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username is incorrect"})
		return
	}

	// Validate the password
	err = utils.ValidatePassword(adminDB.Password, adminRequest.Password)
	if err != nil {
		logger.Error().Msg("Password is incorrect")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect"})
		return
	}

	// Sign a JWT token
	tokenString, err := utils.SignJWT(adminDB.Username)
	if err != nil {
		logger.Error().Msg("Failed to generate JWT token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	// Return the token as response and set the JWT token in the Authorization header
	logger.Info().Msg("Admin successfully logged in")
	c.Header("Authorization", "Bearer "+tokenString)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
