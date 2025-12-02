package main

import (
    "log"
    
    "github.com/gofiber/fiber/v3"
)

func main() {
    // Создаем новое Fiber приложение
    app := fiber.New()

    // Определяем маршрут для GET запроса
    app.Get("/", func(c fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    // Запускаем сервер на порту 3000
    log.Fatal(app.Listen(":3000"))
}