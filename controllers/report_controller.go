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

var reportCollection *mongo.Collection = configs.GetCollection(configs.DB, "reports")

// var validate = validator.New()

func CreateReport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var report models.Report
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&report); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&report); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newUser := models.Report{
		Id:           primitive.NewObjectID(),
		Name:         report.Name,
		Category:     report.Category,
		Url:          report.Url,
		Description:  report.Description,
		Summary:      report.Summary,
		Segmentation: report.Segmentation,
		TOC:          report.TOC,
		Methodology:  report.Methodology,
		Price:        report.Price,
		CreatedAt:    report.CreatedAt,
		UpdatedAt:    report.UpdatedAt,
		Status:       report.Status,
	}

	result, err := reportCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.APIResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAReport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	reportId := c.Params("reportId")
	var report models.Report
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(reportId)

	err := reportCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&report)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": report}})
}

func EditAReport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	reportId := c.Params("reportId")
	var report models.Report
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(reportId)

	//validate the request body
	if err := c.BodyParser(&report); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&report); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	// update := bson.M{"name": report.Name, "location": report.Location, "title": report.Title}
	update := bson.M{
		"name":         report.Name,
		"category":     report.Category,
		"url":          report.Url,
		"description":  report.Description,
		"summary":      report.Summary,
		"segmentation": report.Segmentation,
		"toc":          report.TOC,
		"methodology":  report.Methodology,
		"price":        report.Price,
		"createdAt":    report.CreatedAt,
		"updatedAt":    report.UpdatedAt,
		"status":       report.Status,
	}

	result, err := reportCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated user details
	var updatedUser models.Report
	if result.MatchedCount == 1 {
		err := reportCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

func DeleteAReport(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	reportId := c.Params("reportId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(reportId)

	result, err := reportCollection.DeleteOne(ctx, bson.M{"id": objId})
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

func GetAllReports(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.Report
	defer cancel()

	results, err := reportCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.Report
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		responses.APIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}
