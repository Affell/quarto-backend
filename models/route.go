package models

import "github.com/labstack/echo/v4"

type Route struct {
	Path        string
	Handler     echo.HandlerFunc
	Method      string
	Middlewares []echo.MiddlewareFunc
}
