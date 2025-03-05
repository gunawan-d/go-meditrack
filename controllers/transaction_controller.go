package controllers

import (
	"net/http"
	"strconv"

	"meditrack/repository"

	"github.com/gin-gonic/gin"
)

// TransactionController struct
type TransactionController struct {
	Repo *repository.TransactionRepository
}

func NewTransactionController(repo *repository.TransactionRepository) *TransactionController {
	return &TransactionController{Repo: repo}
}

// CreateTransaction handles medicine purchase by patient
func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var req struct {
		PatientID      int `json:"patient_id"`
		PrescriptionID int `json:"prescription_id"`
		MedicineID     int `json:"medicine_id"`
	}

	// Bind JSON ke struct request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Cek stok obat sebelum transaksi
	medicine, err := tc.Repo.GetMedicineByID(req.MedicineID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Medicine not found"})
		return
	}

	if medicine.Stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// Get Price from medicine
	totalPrice := medicine.Price

	// Create transaction status `pending`
	newID, err := tc.Repo.CreateTransaction(req.PatientID, req.PrescriptionID, req.MedicineID, totalPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        "Transaction created",
		"transaction_id": newID,
		"status":         "pending",
	})
}

// GetTransaction retrieves transaction details
func (tc *TransactionController) GetTransaction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := tc.Repo.GetTransactionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (tc *TransactionController) ProcessPayment(c *gin.Context) {
	var req struct {
		TransactionID int `json:"transaction_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Call if not completed
	err := tc.Repo.CompleteTransactionOnce(req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment processed successfully", "transaction_id": req.TransactionID})
}
