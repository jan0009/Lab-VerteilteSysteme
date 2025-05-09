package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jan0009/Lab-VerteilteSysteme/models"
	"github.com/jan0009/Lab-VerteilteSysteme/storage"
	"gorm.io/gorm"
)

type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateItem(context *fiber.Ctx) error {
	item := Item{}

	err := context.BodyParser(&item)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&item).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create item"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "item has been added"})
	return nil
}

func (r *Repository) DeleteItem(context *fiber.Ctx) error {
	itemModel := models.Items{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	result := r.DB.Delete(&itemModel, id)
	if result.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete item",
		})
		return result.Error
	}

	// Check if anything was actually deleted
	if result.RowsAffected == 0 {
		context.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "item not found",
		})
		return nil
	}

	context.SendStatus(fiber.StatusNoContent)

	return nil
}

func (r *Repository) GetItems(context *fiber.Ctx) error {
	itemModels := &[]models.Items{}

	err := r.DB.Find(itemModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get items"})
		return err
	}

	context.Status(http.StatusOK).JSON(itemModels)

	return nil
}

func (r *Repository) GetItemByID(context *fiber.Ctx) error {

	id := context.Params("id")
	itemModel := models.Items{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&itemModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.Status(http.StatusNotFound).JSON(
			&fiber.Map{"message": "item not found"})
		return nil
	}
	if err != nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "could not get the item"})
		return nil
	}

	context.Status(http.StatusOK).JSON(itemModel)

	return nil
}

func (r *Repository) UpdateItem(context *fiber.Ctx) error {
	item := Item{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	itemID, parseErr := strconv.Atoi(id)
	if parseErr != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "invalid item id",
		})
		return parseErr
	}

	err := context.BodyParser(&item)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"},
		)
		return err
	}

	existingItem := models.Items{}

	err = r.DB.First(&existingItem, itemID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		context.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "item with id not found",
		})
		return nil
	}

	if err != nil {
		context.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "item not found",
		})
		return err
	}

	existingItem.Name = &item.Name
	existingItem.Quantity = &item.Quantity

	saveErr := r.DB.Save(&existingItem).Error
	if saveErr != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not update item",
		})
		return saveErr
	}

	context.Status(http.StatusOK).JSON(existingItem)

	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/items")
	api.Post("/", r.CreateItem)
	api.Delete("/:id", r.DeleteItem)
	api.Put("/:id", r.UpdateItem)
	api.Get("/:id", r.GetItemByID)
	api.Get("/", r.GetItems)
}

func main() {

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173", // oder "*" für alle
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	r.SetupRoutes(app)
	app.Listen(":8080")
}
