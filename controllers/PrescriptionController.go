package controllers

import (
	"database/sql"
	"log"
	"meditrack/middleware"
	"meditrack/repository"
	"meditrack/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrescriptionRepository struct {
	DB *sql.DB
}

func NewPrescriptionRepository(db *sql.DB) *PrescriptionRepository {
	return &PrescriptionRepository{DB: db}
}

// Create recipe Doctor

func CreatePrescriptionHandler(repo repository.PrescriptionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		userClaims, ok := claims.(*middleware.JWTClaims)
		if !ok || userClaims.Role != "doctor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only doctors can process payments"})
			return
		}
		var req struct {
			DoctorID     int    `json:"doctor_id" binding:"required"`
			PatientID    int    `json:"patient_id" binding:"required"`
			Medicine     int    `json:"medicine_id" binding:"required"`
			Dosage       string `json:"dosage" binding:"required"`
			Instructions string `json:"instructions" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		prescription := structs.Prescription{
			DoctorID:     req.DoctorID,
			PatientID:    req.PatientID,
			Medicine:     req.Medicine,
			Dosage:       req.Dosage,
			Instructions: req.Instructions,
			Status:       "pending",
		}

		err := repo.CreatePrescription(prescription)
		if err != nil {
			log.Println("Error creating prescription: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Prescription created"})
	}
}

// Patient Check Receipt
func GetPrescriptionHandler(repo repository.PrescriptionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		patientID, err := strconv.Atoi(c.Param("patient_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
			return
		}

		prescriptions, err := repo.GetPrescriptionByPatientID(patientID)
		if err != nil {
			log.Println("Error getting prescriptions: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, prescriptions)
	}
}

// Doctor reject Receipt
func RejectPrescriptionHandler(repo repository.PrescriptionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get data from JWT
		claims, _ := c.Get("claims")
		userClaims, ok := claims.(*middleware.JWTClaims)
		if !ok || userClaims.Role != "doctor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only doctor can reject prescriptions"})
			return
		}

		// Ambil ID resep dari parameter URL
		param := c.Param("prescription_id")
		log.Println("DEBUG: Prescription ID from URL:", param)

		prescriptionID, err := strconv.Atoi(param)
		if err != nil {
			log.Println("ERROR: Invalid prescription ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prescription ID"})
			return
		}

		// Update status prescription ke "rejected"
		err = repo.UpdatePrescriptionStatus(prescriptionID, "rejected")
		if err != nil {
			log.Println("Error rejecting prescription: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Prescription rejected"})
	}
}

func UpdatePrescriptionStatusHandler(repo repository.PrescriptionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get data from JWT
		claims, _ := c.Get("claims")
		userClaims, ok := claims.(*middleware.JWTClaims)
		if !ok || userClaims.Role != "doctor" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only doctor can update prescriptions"})
			return
		}

		// Ambil ID resep dari parameter URL
		param := c.Param("prescription_id")
		log.Println("DEBUG: Prescription ID from URL:", param)

		prescriptionID, err := strconv.Atoi(param)
		if err != nil {
			log.Println("ERROR: Invalid prescription ID:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid prescription ID"})
			return
		}

		// Ambil status dari body JSON
		var request struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if request.Status != "approved" && request.Status != "rejected" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}

		// Update status prescription
		err = repo.UpdatePrescriptionStatus(prescriptionID, request.Status)
		if err != nil {
			log.Println("Error updating prescription status: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Prescription " + request.Status})
	}
}
