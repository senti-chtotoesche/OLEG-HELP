package jwt

import (
	"time"
)

var JwtKey = []byte("your_secret_key")

func generateJWT(user User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     expirationTime.Unix(),
	}
	token := Jwt.NewWithClaims(Jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
