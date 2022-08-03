package authentication

import (
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := []byte("AllYourBase")
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}
func GenerateJwt(email string, idUser int) (string, error) {
	mySigningKey := []byte("AllYourBase")

	type MyCustomClaims struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		jwt.StandardClaims
	}

	claims := MyCustomClaims{
		idUser,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidJwt(tokenRecive string) (string, error) {
	mySigningKey := []byte("AllYourBase")
	formatToken := strings.Split(tokenRecive, " ")
	type MyCustomClaims struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		jwt.StandardClaims
	}
	token, err := jwt.ParseWithClaims(formatToken[1], &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		println(err.Error())
		return "nil", nil
	}

	if _, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return "Jwt válido", nil
	} else {
		return "", err
	}
}

func ExtractIdOfJwt(tokenRecive string) (int, error) {
	mySigningKey := []byte("AllYourBase")
	formatToken := strings.Split(tokenRecive, " ")
	type MyCustomClaims struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		jwt.StandardClaims
	}
	token, err := jwt.ParseWithClaims(formatToken[1], &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})

	if err != nil {
		println(err.Error())
		return 0, nil
	}

	if _, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return 1, nil
	} else {
		return 0, err
	}
}
