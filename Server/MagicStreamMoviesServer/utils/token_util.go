package utils

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/mnxy-0x/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// for Token creation
type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")
var SECRET_REFRESH_KEY = os.Getenv("SECRET_REFRESH_KEY")

var userCollection = database.OpenCollection("users")

func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {

	// if No validation, and env var is missing, the secret will be an empty string, leading to insecure tokens!
	if SECRET_KEY == "" || SECRET_REFRESH_KEY == "" {
		log.Fatal("SECRET_KEY or SECRET_REFRESH_KEY environment variable is not set")
	}

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(SECRET_REFRESH_KEY))

	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func UpdateAllTokens(userId, token, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) // 2nd_arg will formated by 1st_arg
	if err != nil {
		return err
	}

	updateData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"update_at":     updateAt,
		},
	}
	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData) // (userId) is provided as fn_arg
	if err != nil {
		return err
	}
	return nil
}

func GetAccessToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is required")
	}

	// Check if it actually starts with "Bearer "
	// this approach Won't crash if someone sends "Bearer" (without space) or other malformed headers
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("Authorization header must be Bearer token")
	}

	// authHeader[7:]
	tokenString := authHeader[len(bearerPrefix):]
	if tokenString == "" {
		return "", errors.New("Bearer token is required")
	}

	return tokenString, nil
}
