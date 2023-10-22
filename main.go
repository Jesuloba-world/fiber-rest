package main

import (
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Jesuloba-world/fiber-rest/book"
	"github.com/Jesuloba-world/fiber-rest/database"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/book", book.GetBooks)
	app.Get("/api/v1/book/:id", book.GetBook)
	app.Post("api/v1/book", book.Newbook)
	app.Delete("api/v1/book/:id", book.DeleteBook)
}

func SetupDatabase() (func(), error) {
	// Open a connection to the database
	db, err := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	slog.Info("Database successfully connected")

	// Assign the database connection to the global variable
	database.DB = db

	// Migrate the database schema
	err = database.DB.AutoMigrate(&book.Book{})
	if err != nil {
		return nil, err
	}

	slog.Info("Database migrated")

	closeFunc := func() {
		dbSQL, err := database.DB.DB()
		if err != nil {
			log.Fatal(err)
		}
		dbSQL.Close()
	}

	return closeFunc, nil
}

func main() {
	appConfig := fiber.Config{
		EnablePrintRoutes: true,
		Immutable:         true,
		CaseSensitive:     true,
	}
	app := fiber.New(appConfig)

	slog.Info("Starting app")

	// Set up the database connection
	closeDB, err := SetupDatabase()
	if err != nil {
		panic(err)
	}

	// Close the database connection when the app shuts down
	defer closeDB()

	setupRoutes(app)

	app.Listen(":3000")
}
