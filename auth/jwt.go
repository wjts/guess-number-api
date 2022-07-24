package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TODO: env
var jwtKey = []byte("mu!gK2$M]D*%Ivvm<B[LTKN9\"~GugX+,C%8O;}~498")

type JWTClaim struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(email string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)

	if !ok {
		err = errors.New("claims error")
		return
	}

	if !token.Valid {
		err = errors.New("token is not valid")
		return
	}

	return
}
