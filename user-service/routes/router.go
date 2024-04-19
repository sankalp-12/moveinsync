package routes

import (
	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/user-service/controllers"
	"github.com/sankalp-12/moveinsync/user-service/middlewares"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(users *mongo.Collection, logger zerolog.Logger) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	r.Use(p.Instrument())

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		{
			user.POST("/create", func(c *gin.Context) {
				controllers.Create(c, users, logger)
			})
			user.POST("/login", func(c *gin.Context) {
				controllers.Login(c, users, logger)
			})
			user.POST("/booktrip", middlewares.Validate(logger), func(c *gin.Context) {
				controllers.BookTrip(c, users, logger)
			})
			user.POST("/displaynearbycabs", middlewares.Validate(logger), func(c *gin.Context) {
				controllers.DisplayNearbyCabs(c, users, logger)
			})
		}
	}

	return r
}
