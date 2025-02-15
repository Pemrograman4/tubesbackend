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
	authRoutes.PUT("/users/:id/status", authCtrl.UpdateUserStatus)
	authRoutes.GET("/users", authCtrl.GetAllUsers)
}

	// Course routes
	courseCtrl := controllers.CourseController{DB: db}
	courseRoutes := router.Group("/courses")
	{
		courseRoutes.POST("", courseCtrl.CreateCourse)        // Tambah kursus baru
		courseRoutes.GET("", courseCtrl.GetCourses)           // Dapatkan semua kursus
		courseRoutes.GET("/:id", courseCtrl.FindCourseById)   // Cari kursus berdasarkan ID
		courseRoutes.PUT("/:id", courseCtrl.UpdateCourseById) // Perbarui kursus berdasarkan ID
		courseRoutes.DELETE("/:id", courseCtrl.DeleteCourse)  // Hapus kursus berdasarkan ID

		// Route untuk mendapatkan ID kursus berikutnya
		courseRoutes.GET("/next-id", courseCtrl.GetNextCourseId)
	}

	// Siswa routes
	siswaCtrl := controllers.SiswaController{DB: db}
	siswaRoutes := router.Group("/siswa")
	{
		siswaRoutes.POST("", siswaCtrl.CreateSiswa)
		siswaRoutes.GET("", siswaCtrl.GetSiswa)
		siswaRoutes.GET("/:id", siswaCtrl.GetSiswaByID)
		siswaRoutes.PUT("/:id", siswaCtrl.UpdateSiswa)
		siswaRoutes.DELETE("/:id", siswaCtrl.DeleteSiswa)

	}
	// Guru routes
	guruCtrl := controllers.GuruController{DB: db}
	guruRoutes := router.Group("/gurus")
	{
		guruRoutes.GET("", guruCtrl.GetAllGuru)
		guruRoutes.POST("", guruCtrl.CreateGuru)
		guruRoutes.GET("/:id", guruCtrl.GetGuruByID)
		guruRoutes.PUT("/:id", guruCtrl.UpdateGuru)
		guruRoutes.DELETE("/:id", guruCtrl.DeleteGuru)
	}

	// Tagihan routes
	tagihanCtrl := controllers.TagihanController{DB: db}
	tagihanRoutes := router.Group("/tagihan")
	{
		tagihanRoutes.GET("", tagihanCtrl.GetTagihan)
		tagihanRoutes.GET("/:id", tagihanCtrl.GetTagihanByID)
		tagihanRoutes.POST("", tagihanCtrl.CreateTagihan)
		tagihanRoutes.PUT("/:id", tagihanCtrl.UpdateTagihan)
		tagihanRoutes.DELETE("/:id", tagihanCtrl.DeleteTagihan)
		// Route untuk membayar tagihan
		tagihanRoutes.PUT("/:id/bayar", tagihanCtrl.BayarTagihan)
		tagihanRoutes.GET("/user", tagihanCtrl.GetTagihanByUser)
	}
	return router
}
