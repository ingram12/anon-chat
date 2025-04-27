package main

import (
	"anon-chat-backend/internal/api"
	"anon-chat-backend/internal/config"
	"anon-chat-backend/internal/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	configuration := config.NewConfig()
	userService := handlers.NewUserService(configuration)

	e := echo.New()
	api.RegisterHandlers(e, userService)

	// Добавляем роут для отдачи swagger.json
	e.GET("/swagger.json", func(c echo.Context) error {
		swagger, err := api.GetSwagger()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, swagger)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
