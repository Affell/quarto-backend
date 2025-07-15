package userHandler

import (
	"net/http"
	"quarto/models/user"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUsers récupère la liste des utilisateurs avec pagination
// @Summary Get users list
// @Description Get paginated list of users
// @Tags users
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} user.UserPaginationResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users [get]
func (uh *UserHandler) GetUsers(c echo.Context) error {
	userToken, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	// Paramètres de pagination
	page := 1
	pageSize := 20

	if pageParam := c.QueryParam("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeParam := c.QueryParam("page_size"); pageSizeParam != "" {
		if ps, err := strconv.Atoi(pageSizeParam); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	users, total, err := user.GetUsersPaginated(page, pageSize)
	if err != nil {
		log.Error("Erreur lors de la récupération des utilisateurs", "error", err, "user", userToken.User.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Erreur lors de la récupération des utilisateurs")
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response := user.UserPaginationResponse{
		Users:      users,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}

	return c.JSON(http.StatusOK, response)
}

// GetUser récupère un utilisateur par son ID
// @Summary Get user by ID
// @Description Get user information by ID
// @Tags users
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Param id path int true "User ID"
// @Success 200 {object} user.UserPublic
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (uh *UserHandler) GetUser(c echo.Context) error {
	_, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	userIDParam := c.Param("id")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ID utilisateur invalide")
	}

	userData, err := user.GetUserPublicByID(userID)
	if err != nil {
		if err.Error() == "utilisateur non trouvé" {
			return echo.NewHTTPError(http.StatusNotFound, "Utilisateur non trouvé")
		}
		log.Error("Erreur lors de la récupération de l'utilisateur", "error", err, "requested_user", userID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Erreur lors de la récupération de l'utilisateur")
	}

	return c.JSON(http.StatusOK, userData)
}
