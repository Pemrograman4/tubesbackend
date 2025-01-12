package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/organisasi/tubesbackend/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(db *mongo.Database) *gin.Engine {
	router := gin.Default()

	// Tambahkan middleware CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5504", "http://localhost:5504"}, // Domain frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},        // Metode HTTP yang diperbolehkan
		AllowHeaders:     []string{"Content-Type", "Authorization"},                  // Header yang diizinkan
		ExposeHeaders:    []string{"Content-Length"},                                 // Header yang diizinkan untuk diakses di frontend
		AllowCredentials: true,
	}))

	// Auth routes
	authCtrl := controllers.AuthController{DB: db}
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authCtrl.Register)
		authRoutes.POST("/login", authCtrl.Login)
	}

	// Course routes
	courseCtrl := controllers.CourseController{DB: db}
	courseRoutes := router.Group("/courses")
	{
		courseRoutes.POST("", courseCtrl.CreateCourse)
    courseRoutes.GET("", courseCtrl.GetCourses)
    courseRoutes.PUT("/:id", courseCtrl.UpdateCourse)
    courseRoutes.DELETE("/:id", courseCtrl.DeleteCourse)

    // Tambahkan route baru untuk mendapatkan ID kursus terbaru
    courseRoutes.GET("/latest-id", courseCtrl.GetLatestCourseId)
}
	// Siswa routes
	siswaCtrl := controllers.SiswaController{DB: db}
	siswaRoutes := router.Group("/siswa")
	{
		siswaRoutes.POST("", siswaCtrl.CreateSiswa)
		siswaRoutes.GET("", siswaCtrl.GetSiswa)
		siswaRoutes.GET("/:id", siswaCtrl.GetSiswaByID)
		siswaRoutes.PUT("/:id", siswaCtrl.UpdateSiswa)
		
	return router
}
}