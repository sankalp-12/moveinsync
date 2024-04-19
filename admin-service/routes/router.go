package routes

import (
	"github.com/Depado/ginprom"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/sankalp-12/moveinsync/admin-service/controllers"
	"github.com/sankalp-12/moveinsync/admin-service/middlewares"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(admins *mongo.Collection, cabs *mongo.Collection, logger zerolog.Logger) *gin.Engine {
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
		admin := v1.Group("/admin")
		{
			admin.POST("/create", func(c *gin.Context) {
				controllers.Create(c, admins, logger)
			})
			admin.POST("/login", func(c *gin.Context) {
				controllers.Login(c, admins, logger)
			})
			admin.POST("/addcabs", middlewares.Validate(logger), func(c *gin.Context) {
				controllers.AddCabs(c, cabs, logger)
			})
		}
		cab := v1.Group("/cab")
		{
			cab.GET("/available", func(c *gin.Context) {
				controllers.SuggestAvailableCabs(c, cabs, logger)
			})
			cab.GET("/busy", func(c *gin.Context) {
				controllers.SuggestBusyCabs(c, cabs, logger)
			})
		}
	}

	return r
}
