package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type JwtClaim struct {
	UserID int
	jwt.StandardClaims
}

func (j *JwtWrapper) GenerateToken(UserID int) (signedToken string, err error) {
	claims := &JwtClaim{
		UserID: UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return &JwtClaim{}, err
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		return &JwtClaim{}, errors.New("Couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return &JwtClaim{}, errors.New("JWT expirou")
	}
	return claims, nil
}
