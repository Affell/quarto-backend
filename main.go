package main

import (
	"embed"
	"fmt"
	"os"
	"quarto/config"
	"quarto/handlers"

	_ "quarto/docs"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//go:embed public/*
var Folder embed.FS

// @title Quarto API
// @version 0.0.1
// @description This is the Quarto API documentation.
// @termsOfService https://quarto.fr/terms/

// @contact.name API Support
// @contact.url https://quarto.fr/support
// @contact.email support@quarto.fr
func main() {

	// Initialize echo
	api := echo.New()
	api.HideBanner = true
	api.HTTPErrorHandler = handlers.OnError

	allowedOrigins := []string{
		"https://quarto.affell.fr",
	}
	if log.GetLevel() == log.DebugLevel {
		allowedOrigins = append(allowedOrigins, "*")
	}

	crs := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Access-Control-Allow-Headers", "Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Request-Method", "Access-Control-Request-Headers", handlers.TokenKeyName},
		AllowMethods:     []string{echo.POST, echo.GET, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	})
	api.Use(crs, middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${uri} ${latency_human} ${bytes_in} ${bytes_out} ${remote_ip}` + "\n",
		Output: os.Stdout,
	}))

	api.Use(middleware.BodyLimit(config.Config.BodySizeLimit), middleware.Gzip(), handlers.AuthMiddleware)

	// Register api routes
	for _, handler := range handlers.All() {
		api.Add(handler.Method, handler.Path, handler.Handler, handler.Middlewares...)
	}

	// Swagger
	if log.GetLevel() == log.DebugLevel {
		api.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	// Start public API
	err := api.Start(fmt.Sprintf(":%s", config.Config.ListenPort))
	if err != nil {
		log.Fatal("Public API handler stopped", "error", err)
	}
}
