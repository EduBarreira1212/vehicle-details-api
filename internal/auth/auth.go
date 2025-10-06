package auth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/EduBarreira1212/vehicle-details-api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CreateToken(userID uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()
	permissions["userId"] = userID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	return token.SignedString([]byte(config.LoadConfig().SecretKey))
}

func ValidateToken(c *gin.Context) error {
	tokenString := extractToken(c)
	if tokenString == "" {
		return errors.New("missing bearer token")
	}

	token, err := jwt.Parse(tokenString, returnVerifyKey)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid token")
}

func ExtractUserID(c *gin.Context) (uint64, error) {
	tokenString := extractToken(c)
	if tokenString == "" {
		return 0, errors.New("missing bearer token")
	}

	token, err := jwt.Parse(tokenString, returnVerifyKey)
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["userId"]), 10, 64)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

func extractToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")

	parts := strings.SplitN(auth, " ", 2)

	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}

	return ""
}

func returnVerifyKey(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	secret := config.LoadConfig().SecretKey

	return []byte(secret), nil
}

func GetUserIDFromContext(c *gin.Context) (uint64, bool) {
	v, ok := c.Get("userId")

	if !ok {
		return 0, false
	}

	id, ok := v.(uint64)

	return id, ok
}
