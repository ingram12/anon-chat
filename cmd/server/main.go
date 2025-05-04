package main

import (
	"anon-chat/internal/api"
	"anon-chat/internal/config"
	"anon-chat/internal/handlers"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	middlewareValidator "github.com/oapi-codegen/echo-middleware"
)

func main() {
	e := echo.New()

	// Add limit of 128K to the request body
	e.Use(middleware.BodyLimit("128K"))

	// Add route to serve swagger.json
	e.GET("/swagger.json", func(c echo.Context) error {
		swagger, err := api.GetSwagger()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, swagger)
	})

	// Use OAPI middleware for request validation based on the Swagger spec.
	swagger, err := api.GetSwagger()
	if err != nil {
		panic("Failed to load swagger.json")
	}
	e.Use(middlewareValidator.OapiRequestValidatorWithOptions(swagger, &middlewareValidator.Options{
		Skipper: func(c echo.Context) bool {
			return !strings.HasPrefix(c.Path(), "/api/")
		},
	}))

	configuration := config.NewConfig()
	userService := handlers.NewUserService(configuration)
	api.RegisterHandlers(e, userService)

	if os.Getenv("APP_ENV") == "dev" {
		// Прокси dev-сервера Vite через Go
		target, _ := url.Parse("http://localhost:5173")
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Всё кроме /api отдаём через прокси
		e.GET("/*", echo.WrapHandler(proxy))
	} else {
		// В продакшне отдаем собранный фронт
		e.Static("/", "frontend")
		e.GET("/", func(c echo.Context) error {
			return c.File("frontend/index.html")
		})
	}

	e.Logger.Fatal(e.Start(":8080"))
}
