package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"hospital-backend/internal/repository"
)

type Claims struct {
	StaffID    int `json:"staff_id"`
	HospitalID int `json:"hospital_id"`
	jwt.RegisteredClaims
}

func GenerateToken(staff repository.StaffRecord, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("JWT_SECRET is required")
	}

	claims := Claims{
		StaffID:    staff.ID,
		HospitalID: staff.HospitalID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func Middleware(secret string) gin.HandlerFunc {
	return func(context *gin.Context) {
		authorization := context.GetHeader("Authorization")
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			return
		}

		var claims Claims
		token, err := jwt.ParseWithClaims(parts[1], &claims, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid || claims.HospitalID <= 0 || claims.StaffID <= 0 {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		context.Set("staff_id", claims.StaffID)
		context.Set("hospital_id", claims.HospitalID)
		context.Next()
	}
}
