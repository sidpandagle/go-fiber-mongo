package controllers

import (
	"context"
	"fibgo/configs"
	"fibgo/models"
	"fibgo/responses"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var authCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	fmt.Println(email, password)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if passErr != nil {
		// return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": passErr.Error()}})
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  user.Name,
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWTSECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func Restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
