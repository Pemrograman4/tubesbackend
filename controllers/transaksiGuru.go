package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/organisasi/tubesbackend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransaksiGuruController struct {
	DB *mongo.Database
}

func (ctrl *TransaksiGuruController) CreateTransaksiGuru(c *gin.Context) {
	var transaksiInput struct {
		GuruID string  `json:"guru_id"`
		Amount float64 `json:"amount"`
		Notes  string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&transaksiInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	guruID, err := primitive.ObjectIDFromHex(transaksiInput.GuruID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Guru ID"})
		return
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	createdAt := now.Format("02-01-2006 15:04:05 WIB")
	monthYear := now.Format("01-2006") // Format untuk bulan dan tahun

	var existingTransaksi models.TransaksiGuru
	err = ctrl.DB.Collection("transaksi_guru").FindOne(
		context.TODO(),
		bson.M{
			"guru_id": guruID,
			"created_at": bson.M{"$regex": monthYear}, // Cari transaksi dalam bulan & tahun yang sama
		},
	).Decode(&existingTransaksi)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Guru ini sudah memiliki transaksi di bulan ini"})
		return
	}

	var guru struct {
		FullName string `bson:"fullname"`
	}
	err = ctrl.DB.Collection("gurus").FindOne(context.TODO(), bson.M{"_id": guruID}).Decode(&guru)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch guru data"})
		return
	}

	transaksi := models.TransaksiGuru{
		ID:        primitive.NewObjectID(),
		GuruID:    guruID,
		GuruName:  guru.FullName,
		Amount:    transaksiInput.Amount,
		CreatedAt: createdAt,
		Notes:     transaksiInput.Notes,
	}

	_, err = ctrl.DB.Collection("transaksi_guru").InsertOne(context.TODO(), transaksi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, transaksi)
}

// GetAllTransaksiGuru - Mengambil semua transaksi guru
func (ctrl *TransaksiGuruController) GetAllTransaksiGuru(c *gin.Context) {
	cursor, err := ctrl.DB.Collection("transaksi_guru").Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	var results []models.TransaksiGuru
	if err = cursor.All(context.TODO(), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse transactions"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetTransaksiGuruByID - Mengambil transaksi berdasarkan ID
func (ctrl *TransaksiGuruController) GetTransaksiGuruByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var transaksi models.TransaksiGuru
	err = ctrl.DB.Collection("transaksi_guru").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&transaksi)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaksi)
}

// UpdateTransaksiGuru - Memperbarui transaksi berdasarkan ID
func (ctrl *TransaksiGuruController) UpdateTransaksiGuru(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updateData struct {
		Amount float64 `json:"amount"`
		Notes  string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	update := bson.M{
		"amount": updateData.Amount,
		"notes":  updateData.Notes,
	}

	_, err = ctrl.DB.Collection("transaksi_guru").UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction updated successfully"})
}

// DeleteTransaksiGuru - Menghapus transaksi berdasarkan ID
func (ctrl *TransaksiGuruController) DeleteTransaksiGuru(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = ctrl.DB.Collection("transaksi_guru").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

func (ctrl *TransaksiGuruController) GetLaporanGajiGuru(c *gin.Context) {
    month := c.Query("month") // Format dari frontend: "YYYY-MM"
    if month == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Month parameter is required"})
        return
    }

    // Ekstrak bulan dan tahun dari format "YYYY-MM"
    parts := strings.Split(month, "-")
    if len(parts) != 2 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month format"})
        return
    }
    year, monthNum := parts[0], parts[1]

    // Query untuk mengambil transaksi sesuai bulan dan tahun
    filter := bson.M{
        "created_at": bson.M{
            "$regex": primitive.Regex{Pattern: fmt.Sprintf("^\\d{2}-%s-%s", monthNum, year), Options: "i"},
        },
    }

    var transaksi []models.TransaksiGuru
    cursor, err := ctrl.DB.Collection("transaksi_guru").Find(context.TODO(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
        return
    }
    defer cursor.Close(context.TODO())

    if err = cursor.All(context.TODO(), &transaksi); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode transactions"})
        return
    }

    c.JSON(http.StatusOK, transaksi)
}