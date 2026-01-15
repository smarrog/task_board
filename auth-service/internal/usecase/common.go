package usecase

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	do "github.com/smarrog/task-board/auth-service/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func getPwdHash(pwd do.Pwd) (do.PwdHash, error) {
	pwdHashBytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w: %v", do.ErrPwdHashGeneration, err)
	}
	return do.NewPwdHash(string(pwdHashBytes))
}

func getAccessToken(subject string, secret string, ttl time.Duration) (do.AccessToken, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	kn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := kn.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return do.NewAccessToken(raw)
}
