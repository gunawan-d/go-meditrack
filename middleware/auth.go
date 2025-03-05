package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"

)

type JWTClaims struct {
    Username string `json:"username"`
    UserID   int    `json:"user_id"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func init() {
	if os.Getenv("JWT_SECRET") == "" {
		fmt.Println("JWT_SECRET is not set")
	}
}

func extractToken(c *gin.Context) string {
    bearerToken := c.GetHeader("Authorization")
    if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:7]) == "BEARER " {
        return bearerToken[7:]
    }
    return ""
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		// Save claims di context
		c.Set("claims", claims)

		c.Next()
	}
}

// validateToken Verify token JWT
func validateToken(tokenString, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}