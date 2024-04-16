package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/admin-service/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(admins *mongo.Collection, logger zerolog.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/admin")
		{
			user.POST("/create", func(c *gin.Context) {
				controllers.Create(c, admins, logger)
			})
			user.POST("/login", func(c *gin.Context) {
				controllers.Login(c, admins, logger)
			})
		}
	}

	return r
}
