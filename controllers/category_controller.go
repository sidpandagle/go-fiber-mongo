package controllers

import (
	"context"
	"fibgo/configs"
	"fibgo/models"
	"fibgo/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var categoryCollection *mongo.Collection = configs.GetCollection(configs.DB, "categories")

// var validate = validator.New()

func CreateCategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var category models.Category
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&category); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&category); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newCategory := models.Category{
		Id:          primitive.NewObjectID(),
		Name:        category.Name,
		Description: category.Description,
		Icon:        category.Icon,
	}

	result, err := categoryCollection.InsertOne(ctx, newCategory)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.SuccessAPIResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetACategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	categoryId := c.Params("categoryId")
	var category models.Category
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(categoryId)

	err := categoryCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&category)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": category}})
}

func EditACategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	categoryId := c.Params("categoryId")
	var category models.Category
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(categoryId)

	//validate the request body
	if err := c.BodyParser(&category); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&category); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"name":        category.Name,
		"description": category.Description,
		"icon":        category.Icon,
	}

	result, err := categoryCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated user details
	var updatedUser models.Category
	if result.MatchedCount == 1 {
		err := categoryCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteACategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	categoryId := c.Params("categoryId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(categoryId)

	result, err := categoryCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.SuccessAPIResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllCategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.Category
	defer cancel()

	results, err := categoryCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.Category
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}
