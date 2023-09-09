package fiber_play

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func Fiber_Play() {
	log.Println("hello")

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("GET request")
	})

	app.Post("/", func(c *fiber.Ctx) error {
		return c.SendString("POST request")
	})

	app.Get("/json", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name": "Fiber-Golang",
			"date": time.Now().Format(time.RFC3339Nano),
		})
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		err := c.WriteMessage(websocket.TextMessage, []byte("Hello There. This is echo"))
		// Websocket logic
		for {
			mtype, msg, err := c.ReadMessage()
			if err != nil {
				break
			}
			log.Printf("Read: %s", msg)

			echoMsg := fmt.Sprintf("Echo: %s", string(msg))

			err = c.WriteMessage(mtype, []byte(echoMsg))
			if err != nil {
				break
			}
		}
		log.Println("Error:", err)
	}))

	app.Get("/:param", func(c *fiber.Ctx) error {
		return c.SendString("param: " + c.Params("param"))
	})

	log.Fatal(app.Listen(":3000"))
}
