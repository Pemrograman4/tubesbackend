package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/organisasi/tubesbackend/models"
)

type TagihanController struct {
	DB *mongo.Database
}

// GetTagihan mendapatkan daftar tagihan
func (sc *TagihanController) GetTagihan(c *gin.Context) {
	fmt.Println("GetTagihan called") // Tambahkan log
	collection := sc.DB.Collection("tagihans")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tagihan: " + err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var tagihanList []models.Tagihan
	if err = cursor.All(ctx, &tagihanList); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse tagihan data: " + err.Error()})
		return
	}

	if len(tagihanList) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tagihan found"})
		return
	}

	c.JSON(http.StatusOK, tagihanList)
}

func (ctrl *TagihanController) CreateTagihan(c *gin.Context) {
	var tagihanInput struct {
		SiswaID  string `json:"siswa_id"`
		CourseID string `json:"course_id"`
		DueDate  string `json:"due_date"` // Format: "YYYY-MM-DD"
	}

	// Validate input
	if err := c.ShouldBindJSON(&tagihanInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Convert string IDs to MongoDB types
	siswaID, err := primitive.ObjectIDFromHex(tagihanInput.SiswaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiswaID"})
		return
	}

	courseID, err := primitive.ObjectIDFromHex(tagihanInput.CourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CourseID"})
		return
	}

	// Fetch course to get amount
	var course models.Course
	err = ctrl.DB.Collection("courses").FindOne(context.TODO(), bson.M{"_id": courseID}).Decode(&course)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Determine due date
	var dueDateTime time.Time
	if tagihanInput.DueDate == "" {
		// Default: 7 days from now
		dueDateTime = time.Now().AddDate(0, 0, 7)
	} else {
		// Parse provided due date with format "YYYY-MM-DD"
		dueDateTime, err = time.Parse("2006-01-02", tagihanInput.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DueDate format. Use 'YYYY-MM-DD'"})
			return
		}
	}

	// Convert to primitive.DateTime
	dueDate := primitive.NewDateTimeFromTime(dueDateTime)

	// Create a new Tagihan object
	tagihan := models.Tagihan{
		ID:        primitive.NewObjectID(),
		SiswaID:   siswaID,
		CourseID:  courseID,
		Amount:    course.Cost, // Set amount from course cost
		DueDate:   dueDate,
		Paid:      false,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	// Insert into database
	_, err = ctrl.DB.Collection("tagihans").InsertOne(context.TODO(), tagihan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Tagihan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tagihan created successfully",
		"tagihan": tagihan,
	})
}
func (sc *TagihanController) GetTagihanByUser(c *gin.Context) {
    // Ambil userID dari konteks (misalnya dari token JWT)
    userID := c.GetString("userID")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    collection := sc.DB.Collection("tagihans")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Cari semua tagihan berdasarkan userID
    cursor, err := collection.Find(ctx, bson.M{"user_id": userID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tagihans"})
        return
    }
    defer cursor.Close(ctx)

    var tagihans []models.Tagihan
    if err = cursor.All(ctx, &tagihans); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode tagihans"})
        return
    }

    c.JSON(http.StatusOK, tagihans)
}

// GetTagihanByID mengambil data siswa berdasarkan ID
func (sc *TagihanController) GetTagihanByID(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := sc.DB.Collection("tagihans")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var tagihan models.Tagihan
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&tagihan)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tagihan not found"})
		return
	}

	c.JSON(http.StatusOK, tagihan)
}

func (ctrl *TagihanController) BayarTagihan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Perbarui status menjadi Lunas
	collection := ctrl.DB.Collection("tagihans")
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"paid": true}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tagihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan updated to Lunas"})
}

// UpdateTagihan memperbarui data tagihan berdasarkan ID
func (ctrl *TagihanController) UpdateTagihan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var tagihanInput struct {
		SiswaID  string  `json:"siswa_id"`
		CourseID string  `json:"course_id"`
		Amount   float64 `json:"amount"`
		DueDate  string  `json:"due_date"` // Format ISO 8601 (YYYY-MM-DDTHH:MM:SSZ)
	}

	// Validate input
	if err := c.ShouldBindJSON(&tagihanInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	// Convert string IDs to MongoDB types
	siswaID, err := primitive.ObjectIDFromHex(tagihanInput.SiswaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SiswaID"})
		return
	}

	courseID, err := primitive.ObjectIDFromHex(tagihanInput.CourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid CourseID"})
		return
	}

	// Parse due date
	dueDateTime, err := time.Parse(time.RFC3339, tagihanInput.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DueDate format"})
		return
	}

	// Convert to primitive.DateTime
	dueDate := primitive.NewDateTimeFromTime(dueDateTime)

	collection := ctrl.DB.Collection("tagihans")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Hanya memperbarui data yang diperbolehkan
	update := bson.M{
		"$set": bson.M{
			"siswa_id":  siswaID,
			"course_id": courseID,
			"amount":    tagihanInput.Amount,
			"due_date":  dueDate,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tagihan: " + err.Error()})
		return
	}

	// Mengambil status terbaru dari database (jika perlu)
	var tagihan struct {
		Paid bool `bson:"paid"`
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&tagihan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated tagihan"})
		return
	}

	// Tentukan status berdasarkan field Paid
	tagihanStatus := "Belum Lunas"
	if tagihan.Paid {
		tagihanStatus = "Lunas"
	}

	// Response dengan status
	c.JSON(http.StatusOK, gin.H{
		"message": "Tagihan updated successfully",
		"status":  tagihanStatus,
	})
}

// DeleteTagihan menghapus data tagihan berdasarkan ID
func (ctrl *TagihanController) DeleteTagihan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection := ctrl.DB.Collection("tagihans")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tagihan: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan deleted successfully"})
}
