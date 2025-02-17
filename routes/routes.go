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
	// Membuat instance controller dengan database yang sudah terhubung
	courseCtrl := controllers.NewCourseController(db)
	courseUsersCtrl := controllers.NewCourseUsers(db)
	courseRoutes := router.Group("/courses")
	{
		// Kursus management routes
		courseRoutes.POST("", courseCtrl.CreateCourse)           // Tambah kursus baru
		courseRoutes.GET("", courseCtrl.GetCourses)              // Dapatkan semua kursus
		courseRoutes.GET("/:id", courseCtrl.FindCourseById)      // Cari kursus berdasarkan ID
		courseRoutes.PUT("/:id", courseCtrl.UpdateCourseById)    // Perbarui kursus berdasarkan ID
		courseRoutes.DELETE("/:id", courseCtrl.DeleteCourse)     // Hapus kursus berdasarkan ID
		courseRoutes.GET("/next-id", courseCtrl.GetNextCourseId) // Dapatkan ID kursus berikutnya

		// Pendaftaran kursus
		courseRoutes.POST("/register", courseUsersCtrl.RegisterCourse)                // Daftar kursus
		courseRoutes.GET("/registrations", courseUsersCtrl.GetAllCourseRegistrations) // Dapatkan semua pendaftaran kursus

	}

	// Inisialisasi controller dan rute untuk menangani permintaan
	scheduleCtrl := controllers.NewScheduleController(db)

	scheduleRoutes := router.Group("/schedules")
	scheduleRoutes.Use(middlewares.AuthMiddleware(db)) // Proteksi semua route schedule
	{
		scheduleRoutes.POST("", scheduleCtrl.AddSchedule)                    // Menambahkan jadwal baru
		scheduleRoutes.GET("/:courseId", scheduleCtrl.GetScheduleByCourseId) // Mendapatkan jadwal berdasarkan courseId
		scheduleRoutes.GET("", scheduleCtrl.GetAllSchedules)                 // Mendapatkan semua jadwal
		scheduleRoutes.PUT("/:courseId", scheduleCtrl.UpdateSchedule)        // Memperbarui jadwal berdasarkan courseId
		scheduleRoutes.DELETE("/:courseId", scheduleCtrl.DeleteSchedule)     // Menghapus jadwal berdasarkan courseId
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
		siswaRoutes.POST("/create/transaksi", siswaCtrl.CreateTransaksiSiswa)
		siswaRoutes.PUT("/update/transaksi", siswaCtrl.UpdateStatusTransaksi)
		siswaRoutes.GET("/all/transaksi", siswaCtrl.GetAllTransaksiSiswa)
		siswaRoutes.DELETE("/delete/transaksi/:id", siswaCtrl.DeleteTransaksi)
		siswaRoutes.GET("/get/transaksi/:id", siswaCtrl.GetTransaksiByID)

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
		guruRoutes.GET("/status", guruCtrl.GetGuruByStatus) // Get guru by status
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
		tagihanRoutes.GET("/laporan", tagihanCtrl.GetLaporanTagihan)
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
