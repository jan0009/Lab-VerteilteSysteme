package main

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	"strconv"
)

type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

var items = []Item{
	{ID: 1, Name: "Apples", Quantity: 5},
	{ID: 2, Name: "Bread", Quantity: 2},
	{ID: 3, Name: "Milk", Quantity: 1},
}

func getNextId() int {
	maxId := 0
	for _, item := range items {
		if item.ID > maxId {
			maxId = item.ID
		}
	}
	return maxId + 1
}

func getAllItems(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(items)
}

func getItemById(c *fiber.Ctx) error {
	idParam := c.Params("itemId")
	itemId, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid item ID"})
	}

	for _, item := range items {
		if item.ID == itemId {
			return c.Status(fiber.StatusOK).JSON(item)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Item not found"})
}

func createOrUpdateItem(c *fiber.Ctx) error {
	newItem := new(Item)
	if err := c.BodyParser(newItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	if newItem.Name == "" || newItem.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Name is required and quantity must be a number"})
	}

	for i, item := range items {
		if item.Name == newItem.Name {
			items[i].Quantity += newItem.Quantity
			return c.Status(fiber.StatusOK).JSON(items[i])
		}
	}

	newItem.ID = getNextId()
	items = append(items, *newItem)
	return c.Status(fiber.StatusCreated).JSON(newItem)
}

func updateItem(c *fiber.Ctx) error {
	idParam := c.Params("itemId")
	itemId, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid item ID"})
	}

	updatedItem := new(Item)
	if err := c.BodyParser(updatedItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid input"})
	}

	if updatedItem.Name == "" || updatedItem.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Name is required and quantity must be a number"})
	}

	for i, item := range items {
		if item.ID == itemId {
			for j, other := range items {
				if other.Name == updatedItem.Name && other.ID != itemId {
					items[j].Quantity += updatedItem.Quantity
					items = append(items[:i], items[i+1:]...)
					return c.Status(fiber.StatusOK).JSON(items[j])
				}
			}
			items[i] = Item{ID: itemId, Name: updatedItem.Name, Quantity: updatedItem.Quantity}
			return c.Status(fiber.StatusOK).JSON(items[i])
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Item not found"})
}

func deleteItem(c *fiber.Ctx) error {
	idParam := c.Params("itemId")
	itemId, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid item ID"})
	}

	for i, item := range items {
		if item.ID == itemId {
			items = append(items[:i], items[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Item not found"})
}

func main() {
	app := fiber.New()

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:5173", // oder "*" für alle
	// 	AllowMethods: "GET,POST,PUT,DELETE",
	// }))

	app.Get("/items", getAllItems)
	app.Get("/items/:itemId", getItemById)
	app.Post("/items", createOrUpdateItem)
	app.Put("/items/:itemId", updateItem)
	app.Delete("/items/:itemId", deleteItem)

	app.Listen(":8080")
}
