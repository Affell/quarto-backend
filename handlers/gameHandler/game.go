package gameHandler

import (
	"net/http"
	"quarto/handlers/websocketHandler"
	"quarto/models/game"
	"quarto/models/user"
	"quarto/models/websocket"
	"strconv"

	"github.com/labstack/echo/v4"
)

type GameHandler struct{}

func NewGameHandler() *GameHandler {
	return &GameHandler{}
}

// GetGame récupère une partie
// @Summary Get game
// @Description Get a game by ID
// @Tags games
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Param id path string true "Game ID"
// @Success 200 {object} game.Game
// @Failure 404 {object} map[string]string
// @Router /game/{id} [get]
func (gh *GameHandler) GetGame(c echo.Context) error {
	userToken, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	gameID := c.Param("id")
	g, err := game.GetGame(gameID, userToken.User.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, g.ToWeb())
}

// SelectPiece sélectionne une pièce pour le prochain coup
// @Summary Select piece
// @Description Select a piece for the next move
// @Tags games
// @Accept json
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Param id path string true "Game ID"
// @Param request body game.SelectPieceRequest true "Select piece request"
// @Success 200 {object} game.Game
// @Failure 400 {object} map[string]string
// @Router /game/{id}/select-piece [post]
func (gh *GameHandler) SelectPiece(c echo.Context) error {
	userToken, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	gameID := c.Param("id")
	var req game.SelectPieceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Données invalides")
	}

	updatedGame, err := game.SelectPiece(gameID, userToken.User.ID, req.PieceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Notifier tous les joueurs de la partie via WebSocket
	hub := websocketHandler.GetGameHub(gameID)
	if hub != nil {
		message := websocket.WSMessage{
			Type:   "piece_selected",
			GameID: gameID,
			UserID: strconv.FormatInt(userToken.User.ID, 10),
			Data:   updatedGame.ToWeb(),
		}
		hub.BroadcastToGame(gameID, message)
	}

	return c.JSON(http.StatusOK, updatedGame.ToWeb())
}

// PlacePiece place une pièce sur le plateau
// @Summary Place piece
// @Description Place a piece on the board
// @Tags games
// @Accept json
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Param id path string true "Game ID"
// @Param request body game.PlacePieceRequest true "Place piece request"
// @Success 200 {object} game.Game
// @Failure 400 {object} map[string]string
// @Router /game/{id}/place-piece [post]
func (gh *GameHandler) PlacePiece(c echo.Context) error {
	userToken, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	gameID := c.Param("id")
	var req game.PlacePieceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Données invalides")
	}

	updatedGame, err := game.PlacePiece(gameID, userToken.User.ID, req.Position)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Notifier tous les joueurs de la partie via WebSocket
	hub := websocketHandler.GetGameHub(gameID)
	if hub != nil {
		messageType := "piece_placed"
		if updatedGame.Status == "finished" {
			messageType = "game_finished"
		}

		message := websocket.WSMessage{
			Type:   messageType,
			GameID: gameID,
			UserID: strconv.FormatInt(userToken.User.ID, 10),
			Data:   updatedGame.ToWeb(),
		}
		hub.BroadcastToGame(gameID, message)
	}

	return c.JSON(http.StatusOK, updatedGame.ToWeb())
}

// ForfeitGame abandonne une partie
// @Summary Forfeit game
// @Description Forfeit the current game
// @Tags games
// @Param Quarto-Connect-Token header string true "Session token"
// @Param id path string true "Game ID"
// @Success 200 {object} game.Game
// @Failure 400 {object} map[string]string
// @Router /game/{id}/forfeit [post]
func (gh *GameHandler) ForfeitGame(c echo.Context) error {
	userToken, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	gameID := c.Param("id")
	updatedGame, err := game.ForfeitGame(gameID, userToken.User.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Notifier tous les joueurs de la partie via WebSocket
	hub := websocketHandler.GetGameHub(gameID)
	if hub != nil {
		message := websocket.WSMessage{
			Type:   "game_forfeited",
			GameID: gameID,
			UserID: strconv.FormatInt(userToken.User.ID, 10),
			Data:   updatedGame.ToWeb(),
		}
		hub.BroadcastToGame(gameID, message)
	}

	return c.JSON(http.StatusOK, updatedGame.ToWeb())
}

// GetMyGames récupère toutes les parties de l'utilisateur
// @Summary Get my games
// @Description Get all games for the current user
// @Tags games
// @Produce json
// @Param Quarto-Connect-Token header string true "Session token"
// @Param status query string false "Game status filter (active, finished)"
// @Success 200 {object} []game.Game
// @Failure 401 {object} map[string]string
// @Router /game/my [get]
func (gh *GameHandler) GetMyGames(c echo.Context) error {
	userToken, err := user.GetTokenFromRequest(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	status := c.QueryParam("status")
	var games []game.Game

	if status == "active" {
		games, err = game.GetActiveGames(userToken.User.ID)
	} else {
		games, err = game.GetUserGames(userToken.User.ID)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Convertir en format web
	result := make([]map[string]any, len(games))
	for i, g := range games {
		result[i] = g.ToWeb()
	}

	return c.JSON(http.StatusOK, result)
}
