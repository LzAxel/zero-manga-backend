package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/pkg/clock"
)

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

func (j *JWT) GenerateRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(clock.Now().Add(j.refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(clock.Now()),
		Issuer:    j.issuer,
	})

	signedToken, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
