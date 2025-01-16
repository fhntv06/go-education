package main

import (
	models "database/postgresql/models"
	storage "database/postgresql/storage"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Repository struct {
	DB *gorm.DB
}
type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

func LoadEnvParams() {
	// Загружаем переменные из файла .env
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error LoadEnvParams: %v", err)
	}

	fmt.Println("Loaded env params!")
}
func GetEnvParam(key string) string {
	param := os.Getenv(key)

	if param == "" {
		log.Fatalf("Missing required environment variable: %s!", key)
	}

	return param
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}

	err := context.BodyParser(&book)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book has been added"})

	return nil
}
func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(bookModel, id)
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could bot delete book"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book delete successfully"})
	return nil
}
func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	bookModel := &models.Books{}
	id := context.Params("id")

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Where("id = ?", id).First(bookModel).Error

	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get the book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book id fetched successfully", "data": bookModel})
	return nil
}
func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{} // пакет с моделями

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "books fetched successfully", "data": bookModels})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api") // /api/create_books
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_books/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func ConstructorDB() *gorm.DB {
	config := &storage.Config{
		Host:     GetEnvParam("APP_DB_HOST"),
		Port:     GetEnvParam("APP_DB_PORT"),
		User:     GetEnvParam("APP_DB_USER"),
		Password: GetEnvParam("APP_DB_PASSWORD"),
		DBName:   GetEnvParam("APP_DB_DATABASE"),
		SSLMode:  GetEnvParam("APP_DB_SSL_MODE"),
	}

	fmt.Println("Create storage config:", config)

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatalf("Error could not load the DB: %v", err)
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatalf("could not migrate DB")
	}

	return db
}

func CreateAppFiber() {
	app := fiber.New()

	fmt.Println("Create app fiber!")

	r := Repository{
		DB: ConstructorDB(),
	}

	r.SetupRoutes(app)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatalf("Error Listen 8080: %v", err)
	}
}

func main() {
	LoadEnvParams()

	CreateAppFiber()

	fmt.Print("Hello postgresql!")
}
