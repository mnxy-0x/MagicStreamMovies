package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mnxy-0x/MagicStreamMovies/Server/MagicStreamMoviesServer/database"
	"github.com/mnxy-0x/MagicStreamMovies/Server/MagicStreamMoviesServer/models"
	"github.com/mnxy-0x/MagicStreamMovies/Server/MagicStreamMoviesServer/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = database.OpenCollection("users")

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User

		if err := c.ShouldBindJSON(&user); err != nil {  // ShouldBindJSON(Pointer), fill user_inputs into struct
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user inputs"})
			return
		}

		validate := validator.New()

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation Failed", "details": err.Error()})
			return
		}

		hashedPassword, err := HashPassword(user.Password) // 1. hash password before storing it in DB
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email}) // 2. count how the same email duplicated in DB
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check existing user"}) // 3. failed in count operation
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User is already exists"}) // 4. count success, and duplicated email existed
			return
		}

		user.UserID = bson.NewObjectID().Hex()
		user.CreateAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = hashedPassword

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add the user"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var userCreds models.UserLogin

		if err := c.ShouldBindJSON(&userCreds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid input data"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var foundUser models.User

		// step1: check Email is existed in DB(if user already Signned-up before => authorized access)
		// if email existed; fill (foundUser) with the whole data (User model) from DB
		err := userCollection.FindOne(ctx, bson.M{"email":userCreds.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"email or password is incorrect"})
			return
		}

		// step2: compare inserted password with stored one in DB
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(userCreds.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid email or password"})
			return 
		}
		// if code reach this point; USER CREDS ARE VALID, AND HE IS AUTHORIZED TO ACCESS/LOG-IN
		token, refreshToken, err := utils.GenerateAllTokens(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to generate tokens"})
			return
		}

		err = utils.UpdateAllTokens(foundUser.UserID, token, refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to update tokens"})
			return
		}
		// if code reach here; user successfully logged-in, and tokens are generated fine, SUCCESSFULL_RESPONSE
		c.JSON(http.StatusOK, models.UserResponse{
			UserId: foundUser.UserID,
			FirstName: foundUser.FirstName,
			LastName: foundUser.LastName,
			Email: foundUser.Email,
			Role: foundUser.Role,
			Token: foundUser.Token,
			RefreshToken: foundUser.RefreshToken,
			FavouriteGenres: foundUser.FavouriteGenres,
		})

	}
}