package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/pkg/clock"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (j *JWT) GeneratePair(userID uuid.UUID) (TokenPair, error) {
	var tokenPair = TokenPair{}

	var err error
	tokenPair.AccessToken, err = j.GenerateAccessToken(userID)
	if err != nil {
		return tokenPair, err
	}
	tokenPair.RefreshToken, err = j.GenerateRefreshToken(userID)
	if err != nil {
		return tokenPair, err
	}

	return tokenPair, nil
}

func (j *JWT) GenerateAccessToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(clock.Now().Add(j.accessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(clock.Now()),
		Issuer:    j.issuer,
		Subject:   userID.String(),
	})

	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JWT) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(clock.Now().Add(j.refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(clock.Now()),
		Issuer:    j.issuer,
		Subject:   userID.String(),
	})

	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
