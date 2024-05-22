package utils

import (
	"fmt"
	"errors"
	"os"

	"axiata_test/model"

	jwt "github.com/golang-jwt/jwt/v4"
)

type MyClaims struct {
	jwt.StandardClaims
	Username       string `json:"username"`
	Role           string `json:"role"`
}
var APPLICATION_NAME = "Axiata-Test"
var JWT_SIGNING_METHOD = jwt.SigningMethodHS256

func GenerateTokenJWT(user model.PayloadRegister) string {
	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer: APPLICATION_NAME,
		},
		Username:	user.Username,
		Role:		user.Role,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)
	signatureKey := []byte(os.Getenv("JWT_SIGNATURE_KEY"))

	signedToken, err := token.SignedString(signatureKey)
	if err != nil {
		fmt.Println(" error signing token: ", err)
	}

	return signedToken
}


func DecodeTokenJWT(tokenString string) (resp *model.PayloadRegister, err error) {
	jwtKey := []byte(os.Getenv("JWT_SIGNATURE_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, errors.New("error parsing token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		resp = &model.PayloadRegister{
			Username:	claims["username"].(string),
			Role:		claims["role"].(string),
		}
		return resp, nil
	}
	return nil, errors.New("invalid token")
}