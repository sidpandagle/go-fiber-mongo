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

var swatchCollection *mongo.Collection = configs.GetCollection(configs.DB, "palletes")

// var validate = validator.New()

func CreateSwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var swatch models.Swatch
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&swatch); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&swatch); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newSwatch := models.Swatch{
		Id:        primitive.NewObjectID(),
		Name:      swatch.Name,
		Tags:      swatch.Tags,
		Likes:     swatch.Likes,
		CreatedAt: swatch.CreatedAt,
	}

	result, err := swatchCollection.InsertOne(ctx, newSwatch)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.APIResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetASwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	swatchId := c.Params("swatchId")
	var swatch models.Swatch
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(swatchId)

	err := swatchCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&swatch)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": swatch}})
}

func EditASwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	swatchId := c.Params("swatchId")
	var swatch models.Swatch
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(swatchId)

	//validate the request body
	if err := c.BodyParser(&swatch); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&swatch); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"Name":      swatch.Name,
		"Tags":      swatch.Tags,
		"Likes":     swatch.Likes,
		"CreatedAt": swatch.CreatedAt,
	}

	result, err := swatchCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated user details
	var updatedSwatch models.Swatch
	if result.MatchedCount == 1 {
		err := swatchCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedSwatch)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedSwatch}})
}

func DeleteASwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	swatchId := c.Params("swatchId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(swatchId)

	result, err := swatchCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.APIResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Swatch with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllSwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var swatches []models.Swatch
	defer cancel()

	results, err := swatchCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleSwatch models.Swatch
		if err = results.Decode(&singleSwatch); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		swatches = append(swatches, singleSwatch)
	}

	return c.Status(http.StatusOK).JSON(
		responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": swatches}},
	)
}
