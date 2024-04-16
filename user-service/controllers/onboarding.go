package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/user-service/models"
	"github.com/sankalp-12/moveinsync/user-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Create(c *gin.Context, users *mongo.Collection, logger zerolog.Logger) {
	// Get JSON body from request
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error().Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if the username already exists
	var existingUser models.User
	err := users.FindOne(context.TODO(), bson.M{"username": user.Username}).Decode(&existingUser)
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
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		logger.Error().Msg("Unable to hash password")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to hash password"})
		return
	}

	userModel := models.User{
		Username: user.Username,
		Password: hash,
	}

	// Insert user in MongoDB
	_, err = users.InsertOne(context.TODO(), userModel)
	if err != nil {
		logger.Error().Msg("Unable to insert user in database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to insert user in database"})
		return
	}

	// Return success
	logger.Info().Msg("User successfully created")
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func Login(c *gin.Context, users *mongo.Collection, logger zerolog.Logger) {
	// Get JSON body from request
	var userRequest models.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		logger.Error().Msg("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find the user in MongoDB
	var userDB models.User
	err := users.FindOne(context.TODO(), bson.M{"username": userRequest.Username}).Decode(&userDB)
	if err != nil {
		logger.Error().Msg("Username is incorrect")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username is incorrect"})
		return
	}

	// Validate the password
	err = utils.ValidatePassword(userDB.Password, userRequest.Password)
	if err != nil {
		logger.Error().Msg("Password is incorrect")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect"})
		return
	}

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userDB.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	// Generate the token string
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		logger.Error().Msg("Failed to generate JWT token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	// Return the token as response
	logger.Info().Msg("User successfully logged in")
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
