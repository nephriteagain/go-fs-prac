package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Ballsack")
	app := fiber.New()

	todos := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"todos": todos})
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return c.Status(500).JSON(err)
		}
		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo body is required!"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	app.Patch(("/api/todos/:id"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		if len(id) == 0 {
			return c.SendStatus(400)
		}

		for i, todo := range todos {
			if id == fmt.Sprint(todo.ID) {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	})

	app.Delete(("/api/todos/:id"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		if len(id) == 0 {
			return c.SendStatus(400)
		}
		filteredTodos := []Todo{}
		for _, todo := range todos {
			if id != fmt.Sprint(todo.ID) {
				filteredTodos = append(filteredTodos, todo)
			}
		}
		if len(filteredTodos) == len(todos) {
			return c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("Not found %s", id)})
		}
		todos = filteredTodos
		return c.Status(200).JSON(fiber.Map{"message": fmt.Sprintf("Deleted %s", id)})
	})

	log.Fatal((app.Listen(":4000")))
}
