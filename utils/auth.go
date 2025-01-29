package utils

import (
	"fmt"
	"log"
	"os"

	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Auxesia23/CatatanPengeluaran/models"
	"github.com/golang-jwt/jwt/v5"
)


var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":      user.ID,
		"is_superuser": user.Superuser,
		"username":     user.Username,
		"exp":          time.Now().Add(time.Hour * 24).Unix(),
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
