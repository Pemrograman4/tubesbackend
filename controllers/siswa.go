package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/organisasi/tubesbackend/models"
)

type SiswaController struct {
	DB *mongo.Database
}

// CreateSiswa menambahkan data siswa baru
func (sc *SiswaController) CreateSiswa(c *gin.Context) {
	var siswa models.Siswa
	if err := c.ShouldBindJSON(&siswa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if siswa.FullName == "" || siswa.Address == "" || siswa.PhoneNumber == "" || siswa.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	siswa.ID = primitive.NewObjectID()
	collection := sc.DB.Collection("siswa")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, siswa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create siswa: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

// GetSiswa mendapatkan daftar siswa
func (sc *SiswaController) GetSiswa(c *gin.Context) {
	collection := sc.DB.Collection("siswa")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch siswa: " + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var siswaList []models.Siswa
	if err = cursor.All(ctx, &siswaList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse siswa data: " + err.Error()})
		return
	}

	if len(siswaList) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No students found"})
		return
	}

	c.JSON(http.StatusOK, siswaList)
}

// GetSiswaByID mengambil data siswa berdasarkan ID
func (sc *SiswaController) GetSiswaByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := sc.DB.Collection("siswa")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var siswa models.Siswa
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&siswa)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa not found"})
		return
	}

	c.JSON(http.StatusOK, siswa)
}

// UpdateSiswa memperbarui data siswa berdasarkan ID
func (sc *SiswaController) UpdateSiswa(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var siswa models.Siswa
	if err := c.ShouldBindJSON(&siswa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if siswa.FullName == "" || siswa.Address == "" || siswa.PhoneNumber == "" || siswa.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	collection := sc.DB.Collection("siswa")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"fullname":    siswa.FullName,
			"address":     siswa.Address,
			"phonenumber": siswa.PhoneNumber,
			"email":       siswa.Email,
			"status":      siswa.Status,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update siswa: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa updated successfully"})
}

// DeleteSiswa menghapus data siswa berdasarkan ID
func (sc *SiswaController) DeleteSiswa(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := sc.DB.Collection("siswa")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete siswa: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa deleted successfully"})
}
