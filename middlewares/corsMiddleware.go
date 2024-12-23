package middlewares

import (

	"github.com/gin-gonic/gin"
	"net/http"

)

// ApplyCORS untuk Gin
func ApplyCORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")  // Mengizinkan semua asal
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Jika permintaan adalah preflight (OPTIONS), kirim respons kosong
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusOK)
            return
        }

        c.Next()
    }
}
