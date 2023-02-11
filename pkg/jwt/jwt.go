package jwt

import (
	"fmt"
	"strconv"
	"time"

	gojwt "github.com/golang-jwt/jwt/v4"
)

const UserID = "user_id"

type JWT struct {
	Secret                    string
	RefreshTokenLifespanHours time.Duration
	AuthTokenLifespanMinutes  time.Duration
}

// New returns a new JWT instance
func New(secret, refreshTokenLifespanHours, authTokenLifespanMinutes string) *JWT {
	rl, _ := strconv.Atoi(refreshTokenLifespanHours)
	if rl == 0 {
		rl = 24
	}

	al, _ := strconv.Atoi(authTokenLifespanMinutes)
	if al == 0 {
		al = 5
	}

	return &JWT{
		Secret:                    secret,
		RefreshTokenLifespanHours: time.Duration(rl) * time.Hour,
		AuthTokenLifespanMinutes:  time.Duration(al) * time.Minute,
	}
}

// CreateRefreshToken generates new jwt refresh token with the given user id
func (j *JWT) CreateRefreshToken(userID int64) (string, error) {
	claims := gojwt.MapClaims{}
	claims[UserID] = userID
	claims["exp"] = time.Now().Add(j.RefreshTokenLifespanHours).Unix()
	claims["type"] = "refresh"

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// CreateToken generates new auth token with the given user id
func (j *JWT) CreateToken(userID int64) (string, error) {
	claims := gojwt.MapClaims{}
	claims[UserID] = userID
	claims["exp"] = time.Now().Add(j.AuthTokenLifespanMinutes).Unix()
	claims["type"] = "auth"

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ParseToken validates and decodes a given token and returns a Token object
func (j *JWT) ParseToken(tokenString string) (*gojwt.Token, error) {
	token, err := gojwt.Parse(tokenString, func(token *gojwt.Token) (interface{}, error) {
		// validate the signing method
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return &gojwt.Token{}, err
	}
	return token, nil
}

func (j *JWT) GetUserID(tokenString string) (int64, error) {
	token, err := j.ParseToken(tokenString)
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims := token.Claims.(gojwt.MapClaims)

	id, err := strconv.ParseInt(fmt.Sprintf("%v", claims[UserID]), 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// RefreshToken generates a new token based on the refresh token
func (j *JWT) RefreshToken(refreshTokenString string) (string, error) {
	token, err := j.ParseToken(refreshTokenString)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims := token.Claims.(gojwt.MapClaims)
	if claims["type"] != "refresh" {
		return "", fmt.Errorf("invalid token")
	}

	id, err := strconv.ParseInt(fmt.Sprintf("%v", claims[UserID]), 10, 64)
	if err != nil {
		return "", err
	}
	return j.CreateToken(id)
}
