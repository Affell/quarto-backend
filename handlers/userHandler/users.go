package userHandler

import (
	"net/http"
	"strconv"

	"quarto/models/user"

	"github.com/labstack/echo/v4"
)

// GetUsers returns paginated list of users with optional search
func GetUsers(c echo.Context) error {
	// Parse query parameters
	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("pageSize")
	search := c.QueryParam("search")

	// Default values
	page := 1
	pageSize := 10

	// Parse page
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse pageSize
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Get users with pagination
	result, err := user.GetUsersWithPagination(page, pageSize, search)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

// GetUserByID returns a single user by ID
func GetUserByID(c echo.Context) error {
	// Parse user ID
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Get user by ID
	userPublic, err := user.GetUserByID(int64(id))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, userPublic)
}
