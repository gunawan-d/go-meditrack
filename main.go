package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	"meditrack/controllers"
	"meditrack/database"
	"meditrack/middleware"
	"meditrack/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in environment")
	}

	DB, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	// Migration database
	database.DBMigrate(DB.DB)
	fmt.Println("Successfully connected to database!")

	// Initialize repository
	userRepo := repository.NewUserRepository(DB)
	prescriptionRepo := repository.NewPrescriptionRepository(DB)

	medicineRepo := repository.NewMedicineRepository(DB)
	medicineController := controllers.NewMedicineController(medicineRepo)

	categoryRepo := repository.NewCategoryRepository(DB.DB)
	categoryController := controllers.NewCategoryController(categoryRepo)

	NewTransactionRepository := repository.NewTransactionRepository(DB.DB)

	transactionController := controllers.NewTransactionController(NewTransactionRepository)
	paymentController := controllers.NewPaymentController(NewTransactionRepository)

	//initiate controllers

	// r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())

	// Route User Register & login
	r.POST("/api/users/register", controllers.RegisterHandler(*userRepo))
	r.POST("/api/users/login", controllers.LoginHandler(*userRepo))

	// Protected route with Middleware JWT
	authRoutes := r.Group("/api")
	authRoutes.Use(middleware.JWTMiddleware())

	// Route Prescription
	authRoutes.POST("/prescription", controllers.CreatePrescriptionHandler(*prescriptionRepo))
	authRoutes.GET("/prescription/:patient_id", controllers.GetPrescriptionHandler(*prescriptionRepo))
	authRoutes.PUT("/prescription/:prescription_id/status", controllers.UpdatePrescriptionStatusHandler(*prescriptionRepo))

	// Routes Category
	authRoutes.POST("/categories", categoryController.CreateCategoryHandler)
	authRoutes.GET("/categories", categoryController.GetAllCategoriesHandler)

	// Route Medicine
	authRoutes.POST("/medicines", medicineController.CreateMedicineHandler)
	authRoutes.GET("/medicines", medicineController.GetAllMedicinesHandler)
	authRoutes.PUT("/medicines/:id", medicineController.UpdateMedicineStockHandler)
	authRoutes.DELETE("/medicines/:id", medicineController.DeleteMedicineHandler)

	// Routes Transaction
	authRoutes.POST("/transactions", transactionController.CreateTransaction)
	authRoutes.POST("/payments", transactionController.ProcessPayment)
	authRoutes.GET("/payments/:id", paymentController.GetPaymentStatus)

	r.Run(":8080")
}
