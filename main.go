package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		cursor, err := mg.db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		var employees []Employee = make([]Employee, 0)
		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(employees)

	})
	// create employee
	app.Post("/employee", func(c *fiber.Ctx) error {
		collection := mg.db.Collection("employees")
		employee := new(Employee)
		// if any error exist
		if err := c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		employee.ID = ""
		res, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			c.Status(500).SendString(err.Error())
		}
		// filter
		filter := bson.D{{Key: "_id", Value: res.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)

		createdEmployee := &Employee{}
		createdRecord.Decode(createdEmployee)
		return c.Status(201).JSON(createdEmployee)

	})
	// update employee
	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		// params
		id := c.Params("id")
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.SendStatus(400)
		}
		employee := new(Employee)
		if err := c.BodyParser(employee); err != nil {
			return c.SendStatus(400)
		}
		// query
		query := bson.D{{Key: "_id", Value: oid}}
		update := bson.D{{
			Key: "$set", Value: bson.D{{Key: "name", Value: employee.Name}, {Key: "age", Value: employee.Age}, {Key: "salary", Value: employee.Salary}},
		}}
		err = mg.db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		employee.ID = id
		return c.Status(200).JSON(employee)
	})
	// delete employee
	app.Delete("/employee/:id", func(c *fiber.Ctx) error {

	})

}
