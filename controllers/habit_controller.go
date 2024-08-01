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

var habitCollection *mongo.Collection = configs.GetCollection(configs.DB, "habits")

// var validate = validator.New()

func CreateHabit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var habit models.Habit
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&habit); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&habit); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newHabit := models.Habit{
		Id:       primitive.NewObjectID(),
		Activity: habit.Activity,
		UserId:   habit.UserId,
		Status:   habit.Status,
		Date:     habit.Date,
	}

	result, err := habitCollection.InsertOne(ctx, newHabit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.APIResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAHabit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	habitId := c.Params("habitId")
	var habit models.Habit
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(habitId)

	err := habitCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&habit)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": habit}})
}

func EditAHabit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	habitId := c.Params("habitId")
	var habit models.Habit
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(habitId)

	//validate the request body
	if err := c.BodyParser(&habit); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&habit); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"activity": habit.Activity,
		"userId":   habit.UserId,
		"status":   habit.Status,
		"date":     habit.Date,
	}

	result, err := habitCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated user details
	var updatedUser models.Habit
	if result.MatchedCount == 1 {
		err := habitCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteAHabit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	habitId := c.Params("habitId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(habitId)

	result, err := habitCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.APIResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllHabit(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.Habit
	defer cancel()

	results, err := habitCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.Habit
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}
