package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/organisasi/tubesbackend/controllers"
)

func SetupRoutes(db *mongo.Database) *gin.Engine {
	router := gin.Default()

	authCtrl := controllers.AuthController{DB: db}

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authCtrl.Register)
		authRoutes.POST("/login", authCtrl.Login)
	}

	return router
}
