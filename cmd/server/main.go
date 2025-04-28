package main

import (
	"anon-chat-backend/internal/api"
	"anon-chat-backend/internal/config"
	"anon-chat-backend/internal/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	configuration := config.NewConfig()
	userService := handlers.NewUserService(configuration)

	e := echo.New()

	// Add CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

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
