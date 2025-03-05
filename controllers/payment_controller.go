package controllers

import (
	"net/http"
	"strconv"

	"meditrack/middleware"
	"meditrack/repository"

	"github.com/gin-gonic/gin"
)

// PaymentController struct
type PaymentController struct {
	Repo *repository.TransactionRepository
}

// NewPaymentController creates a new instance
func NewPaymentController(repo *repository.TransactionRepository) *PaymentController {
	return &PaymentController{Repo: repo}
}

// ProcessPayment handles payment for a transaction
func (pc *PaymentController) ProcessPayment(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims, ok := claims.(*middleware.JWTClaims)
	if !ok || userClaims.Role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only doctors can process payments"})
		return
	}

	type PaymentRequest struct {
		TransactionID int `json:"transaction_id"`
	}

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Update status transaksi menjadi "paid"
	err := pc.Repo.UpdateTransactionStatus(req.TransactionID, "completed")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}

	// Kurangi stok obat setelah pembayaran berhasil
	err = pc.Repo.ReduceMedicineStock(req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment successful, but failed to update stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful, stock updated"})
}

// GetPaymentStatus retrieves payment status
func (pc *PaymentController) GetPaymentStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	transaction, err := pc.Repo.GetTransactionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": transaction.Status})
}
