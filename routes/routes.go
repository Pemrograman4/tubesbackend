package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/organisasi/tubesbackend/controllers"
	"github.com/organisasi/tubesbackend/middlewares"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(db *mongo.Database) *gin.Engine {
	router := gin.Default()

	// Middleware CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5504", "http://localhost:5504"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Auth routes
	authCtrl := controllers.AuthController{DB: db}
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authCtrl.Register)
		authRoutes.POST("/login", authCtrl.Login)

		// Gunakan middleware untuk melindungi route ini
		authRoutes.Use(middlewares.AuthMiddleware(db))
		authRoutes.GET("/users", authCtrl.GetAllUsers)
		authRoutes.PUT("/users/:id/status", authCtrl.UpdateUserStatus)
	}

	// Course routes
	courseCtrl := controllers.CourseController{DB: db}
	courseRoutes := router.Group("/courses")
	courseRoutes.Use(middlewares.AuthMiddleware(db)) // Proteksi semua route kursus
	{
		courseRoutes.POST("", courseCtrl.CreateCourse)
		courseRoutes.GET("", courseCtrl.GetCourses)
		courseRoutes.GET("/:id", courseCtrl.FindCourseById)
		courseRoutes.PUT("/:id", courseCtrl.UpdateCourseById)
		courseRoutes.DELETE("/:id", courseCtrl.DeleteCourse)
		courseRoutes.GET("/next-id", courseCtrl.GetNextCourseId)
	}

	// Siswa routes
	siswaCtrl := controllers.SiswaController{DB: db}
	siswaRoutes := router.Group("/siswa")
	siswaRoutes.Use(middlewares.AuthMiddleware(db)) // Proteksi semua route siswa
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
	guruRoutes.Use(middlewares.AuthMiddleware(db)) // Proteksi semua route guru
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
	tagihanRoutes.Use(middlewares.AuthMiddleware(db)) // Proteksi semua route tagihan
	{
		tagihanRoutes.GET("", tagihanCtrl.GetTagihan)
		tagihanRoutes.GET("/:id", tagihanCtrl.GetTagihanByID)
		tagihanRoutes.POST("", tagihanCtrl.CreateTagihan)
		tagihanRoutes.PUT("/:id", tagihanCtrl.UpdateTagihan)
		tagihanRoutes.DELETE("/:id", tagihanCtrl.DeleteTagihan)
		tagihanRoutes.PUT("/:id/bayar", tagihanCtrl.BayarTagihan)
		tagihanRoutes.GET("/user", tagihanCtrl.GetTagihanByUser)
	}
// Transaksi Guru Routes
	transaksiGuruCtrl := controllers.TransaksiGuruController{DB: db}
	transaksiRoutes := router.Group("/transaksi-guru")
	// transaksiRoutes.Use(middlewares.AuthMiddleware(db)) // Proteksi dengan autentikasi
	{

		transaksiRoutes.POST("", transaksiGuruCtrl.CreateTransaksiGuru)
		transaksiRoutes.GET("", transaksiGuruCtrl.GetAllTransaksiGuru)
		transaksiRoutes.GET("/:id", transaksiGuruCtrl.GetTransaksiGuruByID)
		transaksiRoutes.PUT("/:id", transaksiGuruCtrl.UpdateTransaksiGuru)
		transaksiRoutes.DELETE("/:id", transaksiGuruCtrl.DeleteTransaksiGuru)
	}

	// Laporan Guru Routes
	laporanRoutes := router.Group("/laporan-guru")
	// laporanRoutes.Use(middlewares.AuthMiddleware(db)) 
	{

		laporanRoutes.GET("/", transaksiGuruCtrl.GenerateLaporan)
	}
	return router
}
