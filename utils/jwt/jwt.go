package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const UserID = "user_id"

type JWT struct {
	Secret                    string
	RefreshTokenLifespanHours time.Duration
	AuthTokenLifespanMinutes  time.Duration
}

func New(secret, refreshTokenLifespanHours, authTokenLifespanMinutes string) *JWT {
	_refreshTokenLifespanHours, _ := strconv.Atoi(refreshTokenLifespanHours)
	if _refreshTokenLifespanHours == 0 {
		_refreshTokenLifespanHours = 24
	}

	_authTokenLifespanMinutes, _ := strconv.Atoi(authTokenLifespanMinutes)
	if _authTokenLifespanMinutes == 0 {
		_authTokenLifespanMinutes = 5
	}
	return &JWT{
		Secret:                    secret,
		RefreshTokenLifespanHours: time.Duration(_refreshTokenLifespanHours) * time.Hour,
		AuthTokenLifespanMinutes:  time.Duration(_authTokenLifespanMinutes) * time.Minute,
	}
}

// CreateRefreshToken generates new jwt refresh token with the given user id
func (j *JWT) CreateRefreshToken(userID int64) (string, error) {
	claims := jwt.MapClaims{}

	claims[UserID] = userID
	claims["exp"] = time.Now().Add(j.RefreshTokenLifespanHours).Unix()
	claims["type"] = "refresh"
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// CreateToken generates new auth token with the given user id
func (j *JWT) CreateToken(userID int64) (string, error) {
	claims := jwt.MapClaims{}

	claims[UserID] = userID
	claims["exp"] = time.Now().Add(j.AuthTokenLifespanMinutes).Unix()
	claims["type"] = "auth"
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken validates and decodes a given token and returns a Token object
func (j *JWT) ParseToken(tokenString string) (*jwt.Token, error) {
	// validate and extract token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %s", token.Method.Alg())
		}
		return []byte(j.Secret), nil
	})
	if err != nil {
		return &jwt.Token{}, err
	}
	return token, nil
}

func (j *JWT) GetUserID(t string) (int64, error) {
	token, err := j.ParseToken(t)
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)

	id, err := strconv.ParseInt(fmt.Sprintf("%v", claims[UserID]), 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// RefreshToken generates a new token based on the refresh token
func (j *JWT) RefreshToken(refreshToken string) (newToken string, err error) {
	token, err := j.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)

	if claims["type"] != "refresh" {
		return "", fmt.Errorf("invalid token")
	}
	id, err := strconv.ParseInt(fmt.Sprintf("%v", claims[UserID]), 10, 64)
	if err != nil {
		return "", err
	}
	return j.CreateToken(id)
}
