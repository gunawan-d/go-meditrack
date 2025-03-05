package controllers

import (
	"net/http"
	"os"
	"time"

	"meditrack/repository"
	"meditrack/structs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var JwtKey []byte

func init() {
	JwtKey = []byte(os.Getenv("JWT_SECRET"))
	if len(JwtKey) == 0 {
		panic("JWT_SECRET is not set in environment")
	}
}

func LoginHandler(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		user, err := userRepo.GetUserByemail(req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		// Generate JWT token
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &struct {
			Email  string `json:"email"`
			UserID int    `json:"user_id"`
			Role   string `json:"role"`
			jwt.StandardClaims
		}{
			Email:  user.Email,
			UserID: user.ID,
			Role:   user.Role,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(JwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

// RegisterHandler
func RegisterHandler(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		emailExists, err := userRepo.IsEmailExists(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check email"})
			return
		}
		if emailExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := structs.User{
			Name:     req.Name,
			Role:     req.Role,
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		err = userRepo.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}
