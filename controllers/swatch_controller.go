package controllers

import (
	"context"
	"fibgo/configs"
	"fibgo/models"
	"fibgo/responses"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var swatchCollection *mongo.Collection = configs.GetCollection(configs.DB, "palletes")

// var validate = validator.New()

func CreateSwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var swatch models.Swatch
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&swatch); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&swatch); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newSwatch := models.Swatch{
		Id:        primitive.NewObjectID(),
		Name:      swatch.Name,
		Tags:      swatch.Tags,
		Likes:     swatch.Likes,
		CreatedAt: time.Now(),
	}

	result, err := swatchCollection.InsertOne(ctx, newSwatch)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.SuccessAPIResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func CreateSwatchList(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var swatches []models.Swatch

	// Validate the request body
	if err := c.BodyParser(&swatches); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Iterate through the array and validate each swatch
	for _, swatch := range swatches {
		if validationErr := validate.Struct(swatch); validationErr != nil {
			return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
		}
	}

	// Create an array to store the new swatches
	var newSwatches []interface{}

	// Iterate through the array and create a new swatch for each
	for _, swatch := range swatches {
		newSwatch := models.Swatch{
			Id:        primitive.NewObjectID(),
			Name:      swatch.Name,
			Tags:      swatch.Tags,
			Likes:     swatch.Likes,
			CreatedAt: time.Now(),
		}
		newSwatches = append(newSwatches, newSwatch)
	}

	// Insert the array of new swatches into the database
	result, err := swatchCollection.InsertMany(ctx, newSwatches)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.SuccessAPIResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetASwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	swatchId := c.Params("swatchId")
	var swatch models.Swatch
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(swatchId)

	err := swatchCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&swatch)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": swatch}})
}

func EditASwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	swatchId := c.Params("swatchId")
	var swatch models.Swatch
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(swatchId)

	//validate the request body
	if err := c.BodyParser(&swatch); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&swatch); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.SuccessAPIResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{
		"Name":      swatch.Name,
		"Tags":      swatch.Tags,
		"Likes":     swatch.Likes,
		"CreatedAt": time.Now(),
	}

	result, err := swatchCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated user details
	var updatedSwatch models.Swatch
	if result.MatchedCount == 1 {
		err := swatchCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedSwatch)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedSwatch}})
}

func IncrementLike(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	swatchId := c.Params("swatchId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(swatchId)

	// Increment the Likes field by 1
	update := bson.M{
		"$inc": bson.M{"Likes": 1},
		"$set": bson.M{"CreatedAt": time.Now()},
	}

	result, err := swatchCollection.UpdateOne(ctx, bson.M{"id": objId}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Get updated swatch details
	var updatedSwatch models.Swatch
	if result.MatchedCount == 1 {
		err := swatchCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedSwatch)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedSwatch}})
}

func DeleteASwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	swatchId := c.Params("swatchId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(swatchId)

	result, err := swatchCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.SuccessAPIResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Swatch with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllSwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var swatches []models.Swatch
	defer cancel()

	results, err := swatchCollection.Find(ctx, bson.M{}, options.Find().SetLimit(10).SetSkip(10))

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleSwatch models.Swatch
		if err = results.Decode(&singleSwatch); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		// Calculate time difference with the current time
		timeDifference := time.Since(singleSwatch.CreatedAt)

		// Add the time difference to the Swatch struct
		singleSwatch.TimeDifference = timeDifference

		swatches = append(swatches, singleSwatch)
	}

	return c.Status(http.StatusOK).JSON(
		responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": swatches}},
	)
}

func GetFilteredSwatch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var swatches []models.Swatch
	defer cancel()

	limitStr := c.Query("limit")
	limit, errl := strconv.ParseInt(limitStr, 10, 64)
	if errl != nil {
		limit = 10
	}

	pageStr := c.Query("page")
	page, errl := strconv.ParseInt(pageStr, 10, 64)
	if errl != nil {
		page = 1
	}

	trendingStr := c.Query("trending")
	trending, errl := strconv.ParseInt(trendingStr, 10, 64)
	if errl != nil {
		trending = 0
	}

	// Get the search query parameter from the request
	searchQuery := c.Query("search")

	// Create a filter based on the search query
	filter := bson.M{}
	fmt.Println("Val:", searchQuery, searchQuery == "")
	if searchQuery != "" {
		tags := strings.Split(searchQuery, ",")
		filter["tags"] = bson.M{"$all": tags}
	}

	findOptions := options.Find().SetLimit(int64(limit)).SetSkip((int64(page) - 1) * int64(limit)).SetSort(bson.D{{"createdAt", 1}})

	if trending == 1 {
		findOptions.SetSort(bson.D{{"likes", -1}})
	}

	results, err := swatchCollection.Find(ctx, filter, findOptions)
	if searchQuery != "" {
		fmt.Println("Filter applied:", filter)
	}
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleSwatch models.Swatch
		if err = results.Decode(&singleSwatch); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.SuccessAPIResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		// Calculate time difference with the current time
		timeDifference := time.Since(singleSwatch.CreatedAt)

		// Add the time difference to the Swatch struct
		singleSwatch.TimeDifference = timeDifference

		swatches = append(swatches, singleSwatch)
	}

	return c.Status(http.StatusOK).JSON(
		responses.SuccessAPIResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": swatches}},
	)
}
