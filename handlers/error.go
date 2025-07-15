package handlers

import (
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

func OnError(err error, c echo.Context) {
	fmtError := Error{
		Code:    http.StatusInternalServerError,
		Message: "Erreur interne au serveur",
	}

	req := c.Request()

	he, ok := err.(*echo.HTTPError)
	if ok {
		switch message := he.Message.(type) {
		case string:
			fmtError.Message = message
		case error:
			fmtError.Message = message.Error()
		}
		fmtError.Code = he.Code

	}

	log.Warnf("[%s - %s] %s", req.Method, req.RequestURI, err)

	_ = c.JSONPretty(fmtError.Code, fmtError, "\t")
}
