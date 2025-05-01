package main

import (
	"anon-chat/internal/api"
	"anon-chat/internal/config"
	"anon-chat/internal/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	middlewareValidator "github.com/oapi-codegen/echo-middleware"
)

func main() {
	e := echo.New()

	// Add limit of 128K to the request body
	e.Use(middleware.BodyLimit("128K"))

	// Add CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	swagger, err := api.GetSwagger()
	if err != nil {
		panic("Failed to load swagger.json")
	}

	// Use OAPI middleware for request validation based on the Swagger spec.
	e.Use(middlewareValidator.OapiRequestValidator(swagger))

	configuration := config.NewConfig()
	userService := handlers.NewUserService(configuration)
	api.RegisterHandlers(e, userService)

	// Add route to serve swagger.json
	e.GET("/swagger.json", func(c echo.Context) error {
		swagger, err := api.GetSwagger()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, swagger)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
