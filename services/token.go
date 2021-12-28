package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"tianchi/dao/cache"
	db "tianchi/dao/mysql"
	"tianchi/models"
)

var privateKey = []byte("privateKey")

type CustomClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func CreateToken(user *models.User) (token string, err error) {
	if user == nil {
		return "", errors.New("user is nil")
	}

	var claim CustomClaims
	claim.UserId = user.Username
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err = at.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	cache.UserToken.Store(token, user.Username)
	db.WriteUserTokenChan <- &models.Token{
		Username: user.Username,
		Token:    token,
	}
	return token, nil
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return privateKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
