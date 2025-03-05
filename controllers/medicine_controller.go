package controllers

import (
	"meditrack/middleware"
	"meditrack/repository"
	"meditrack/structs"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MedicineController struct {
	Repo *repository.MedicineRepository
}

func NewMedicineController(repo *repository.MedicineRepository) *MedicineController {
	return &MedicineController{Repo: repo}
}

// Add Medicine
func (c *MedicineController) CreateMedicineHandler(ctx *gin.Context) {
	var medicine structs.Medicine
	if err := ctx.ShouldBindJSON(&medicine); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := c.Repo.CreateMedicine(medicine)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Medicine created successfully"})
}

// Get All medicines
func (c *MedicineController) GetAllMedicinesHandler(ctx *gin.Context) {
	medicines, err := c.Repo.GetAllMedicines()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, medicines)
}

// Update Stock Medicine
func (c *MedicineController) UpdateMedicineStockHandler(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine id"})
		return
	}

	var updateData struct {
		Stock int `json:"stock" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&updateData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = c.Repo.UpdateMedicineStock(id, updateData.Stock)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}

func (mc *MedicineController) DeleteMedicineHandler(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims, ok := claims.(*middleware.JWTClaims)
	if !ok || userClaims.Role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only doctors can delete medicine"})
		return
	}

	medicineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medicine ID"})
		return
	}

	err = mc.Repo.DeleteMedicine(medicineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete medicine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Medicine deleted successfully"})
}
