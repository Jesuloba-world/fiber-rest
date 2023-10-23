package book

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/Jesuloba-world/fiber-rest/database"
)

// model
type Book struct {
	gorm.Model
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

func GetBooks(c *fiber.Ctx) error {
	db := database.DB
	var books []Book
	db.Find(&books)
	return c.JSON(books)
}

func GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		return c.Status(http.StatusNotFound).SendString("No book found with given id")
	}
	return c.JSON(book)
}
func Newbook(c *fiber.Ctx) error {
	db := database.DB

	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	db.Create(&book)
	return c.JSON(book)
}

func DeleteBook(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		return c.Status(http.StatusNotFound).SendString("No book found with given id")
	}
	db.Delete(&book)
	return c.Status(http.StatusOK).SendString("Book successfully deleted")
}

func UpdateBook(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.DB

	// Parse the updated book details from the request body
	updatedBook := new(Book)
	if err := c.BodyParser(updatedBook); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	// Update the book with the new details
	result := db.Model(&Book{}).Where("id = ?", id).Updates(updatedBook)
	if result.Error != nil {
		return c.Status(http.StatusNotFound).SendString("No book found with given id")
	}

	return c.JSON(updatedBook)
}
