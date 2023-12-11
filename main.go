package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	client *mongo.Client
	db     *mongo.Database
}

// instance
var mg MongoInstance

// db name
const dbName = "Go-Fiber-hrms"

// connection string
const mongoURI = os.Getenv("DB")

// struct
type Employee struct {
	ID     string
	Name   string
	Salary float64
	Age    float64
}

// database connection function
func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	// timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	db := client.Database(dbName)
	if err != nil {
		panic(err)
	}
	mg = MongoInstance{
		client: client,
		db:     db,
	}
	return nil
}

func main() {
	if err := Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	app.Get("/employee", func(c *fiber.Ctx) error {

	})
	app.Post("/employee", func(c *fiber.Ctx) error {

	})
	app.Put("/employee/:id", func(c *fiber.Ctx) error {

	})
	app.Delete("/employee/:id", func(c *fiber.Ctx) error {

	})

}
