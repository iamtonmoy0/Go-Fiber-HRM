package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
var mongoURI = os.Getenv("DB")

// struct
type Employee struct {
	ID     string  `json:"id",omitempty bson:"_id" ,omitempty`
	Name   string  `json:"name" `
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
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
	// get employee
	app.Get("/employee", func(c *fiber.Ctx) error {
    query := bson.D{{}}
   cursor, err:= mg.db.Collection("employees").Find(c.Context(), query)
   if err != nil {
	return c.Status(400).SendString(err.Error())
   }
   var employees []Employees = make([]Employees, 0)
   if err:=cursor.All(c.Context(),&employees);err!=nil{
	return c.Status(500).SendStatus(err.Error())
   }
   return c.JSON(employees)

	}
)
	// create employee
	app.Post("/employee", func(c *fiber.Ctx) error {

	})
	// update employee
	app.Put("/employee/:id", func(c *fiber.Ctx) error {

	})
	// delete employee
	app.Delete("/employee/:id", func(c *fiber.Ctx) error {

	})

}
