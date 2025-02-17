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

	// Convert string IDs to MongoDB ObjectIDs
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

	// Fetch siswa data (ambil nama & email)
	var siswa models.Siswa
	err = ctrl.DB.Collection("siswa").FindOne(context.TODO(), bson.M{"_id": siswaID}).Decode(&siswa)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa not found"})
		return
	}

	// Fetch course data (ambil nama & harga)
	var course models.Course
	err = ctrl.DB.Collection("courses").FindOne(context.TODO(), bson.M{"_id": courseID}).Decode(&course)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Determine due date
	var dueDateTime time.Time
	if tagihanInput.DueDate == "" {
		dueDateTime = time.Now().AddDate(0, 0, 7) // Default: 7 hari dari sekarang
	} else {
		dueDateTime, err = time.Parse("2006-01-02", tagihanInput.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DueDate format. Use 'YYYY-MM-DD'"})
			return
		}
	}
	dueDate := primitive.NewDateTimeFromTime(dueDateTime)

	// Create tagihan dengan data siswa & course
	tagihan := models.Tagihan{
		ID:        primitive.NewObjectID(),
		SiswaID:   siswaID,
		SiswaNama: siswa.FullName,
		SiswaEmail: siswa.Email,
		CourseID:  courseID,
		CourseName: course.Name,
		Amount:    course.Cost,
		DueDate:   dueDate,
		Paid:      false,
		Status:    "Belum Bayar", // Set default status saat tagihan dibuat
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	// Insert into database
	_, err = ctrl.DB.Collection("tagihans").InsertOne(context.TODO(), tagihan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Tagihan"})
		return
	}

	// Response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Tagihan created successfully",
		"tagihan": tagihan,
	})
}

func (sc *TagihanController) GetTagihanByUser(c *gin.Context) {
	// Ambil user_id dari context yang disimpan di middleware
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Konversi userID ke ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Cari user berdasarkan ID
	var user models.User
	err = sc.DB.Collection("users").FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Cari siswa berdasarkan email
	var siswa models.Siswa
	err = sc.DB.Collection("siswa").FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&siswa)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa not found"})
		return
	}

	// Cari semua tagihan berdasarkan email siswa
	collection := sc.DB.Collection("tagihans")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"siswa_email": siswa.Email}) // Gunakan email
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

	// Waktu saat ini
	now := primitive.NewDateTimeFromTime(time.Now())

	// Perbarui status menjadi Lunas, serta set paid_at dan updated_at
	collection := ctrl.DB.Collection("tagihans")
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"paid":      true,
			"status":    "Lunas",
			"paid_at":   now,
			"updated_at": now,
		}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tagihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan updated to Lunas", "paid_at": now, "updated_at": now})
}

func (ctrl *TagihanController) UpdateTagihan(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Parse request body untuk mendapatkan data yang ingin diupdate
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updateFields := bson.M{}

	// Jika due_date ada, konversi string ke DateTime
	if dueDateStr, exists := updateData["due_date"].(string); exists {
		dueDateTime, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DueDate format. Use 'YYYY-MM-DD'"})
			return
		}
		updateFields["due_date"] = primitive.NewDateTimeFromTime(dueDateTime)
	}

	// Jika siswa_id ada, konversi string ke ObjectID dan ambil data siswa
	if siswaIDStr, exists := updateData["siswa_id"].(string); exists {
		siswaID, err := primitive.ObjectIDFromHex(siswaIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid siswa_id"})
			return
		}
		updateFields["siswa_id"] = siswaID

		var siswa models.Siswa
		err = ctrl.DB.Collection("siswa").FindOne(context.TODO(), bson.M{"_id": siswaID}).Decode(&siswa)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Siswa not found"})
			return
		}
		updateFields["siswa_nama"] = siswa.FullName
		updateFields["siswa_email"] = siswa.Email
	}

	// Jika course_id ada, konversi string ke ObjectID dan ambil data course
	if courseIDStr, exists := updateData["course_id"].(string); exists {
		courseID, err := primitive.ObjectIDFromHex(courseIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course_id"})
			return
		}
		updateFields["course_id"] = courseID

		var course models.Course
		err = ctrl.DB.Collection("courses").FindOne(context.TODO(), bson.M{"_id": courseID}).Decode(&course)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		updateFields["course_name"] = course.Name
	}

	// Pastikan updated_at tidak diambil dari input user
	updateFields["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	// Perbarui tagihan berdasarkan data yang diberikan
	collection := ctrl.DB.Collection("tagihans")
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tagihan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tagihan updated successfully"})
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

func (ctrl *TagihanController) GetLaporanTagihan(c *gin.Context) {
    status := c.QueryArray("status") // Mengambil status dari query parameter, bisa kosong
    startDate, _ := time.Parse("2006-01-02", c.Query("start_date"))
    endDate, _ := time.Parse("2006-01-02", c.Query("end_date"))
    filter := bson.M{}

    // Jika ada status yang diterima, lakukan filter berdasarkan status
    if len(status) > 0 {
        filter["status"] = bson.M{"$in": status}
    }

    // Jika ada rentang tanggal, filter berdasarkan tanggal
    if !startDate.IsZero() && !endDate.IsZero() {
        filter["created_at"] = bson.M{"$gte": startDate, "$lte": endDate}
    }

    collection := ctrl.DB.Collection("tagihans")
    cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil laporan"})
        return
    }

    var tagihans []models.Tagihan
    if err = cursor.All(context.TODO(), &tagihans); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca data"})
        return
    }

    c.JSON(http.StatusOK, tagihans)
}