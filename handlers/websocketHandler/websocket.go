package websocketHandler

import (
	"net/http"
	"quarto/models/user"
	"quarto/models/websocket"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type WebSocketHandler struct {
	gameHubs map[string]*websocket.Hub
	mutex    sync.RWMutex
}

var gameHandler *WebSocketHandler

func init() {
	gameHandler = &WebSocketHandler{
		gameHubs: make(map[string]*websocket.Hub),
	}
}

func NewWebSocketHandler() *WebSocketHandler {
	return gameHandler
}

// GetOrCreateGameHub retourne ou crée un hub pour une partie spécifique
func (wsh *WebSocketHandler) GetOrCreateGameHub(gameID string) *websocket.Hub {
	wsh.mutex.Lock()
	defer wsh.mutex.Unlock()

	if hub, exists := wsh.gameHubs[gameID]; exists {
		return hub
	}

	// Créer un nouveau hub pour cette partie
	hub := websocket.NewHub()
	go hub.Run()
	wsh.gameHubs[gameID] = hub
	return hub
}

// HandleWebSocket gère les connexions WebSocket
// @Summary WebSocket connection
// @Description Establish WebSocket connection for real-time communication
// @Tags websocket
// @Param game_id query string true "Game ID (required for game-specific communication)"
// @Param token query string true "Session token"
// @Router /ws [get]
func (wsh *WebSocketHandler) HandleWebSocket(c echo.Context) error {

	strToken := c.QueryParam("token")

	if len(strToken) != 36 {
		return echo.NewHTTPError(http.StatusBadRequest, "paramètre token requis")
	}

	userToken, err := user.GetUserToken(strToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "token invalide")
	}

	if time.Since(userToken.CreatedAt) > user.TOKEN_EXPIRATION {
		user.RevokeUserToken(userToken.TokenID)
		return echo.NewHTTPError(http.StatusBadRequest, "token expiré")
	}

	gameID := c.QueryParam("game_id")

	if gameID == "" {
		return echo.NewHTTPError(400, "game_id requis pour la connexion WebSocket")
	}

	// Connexion pour une partie spécifique
	hub := wsh.GetOrCreateGameHub(gameID)
	return hub.HandleWebSocket(c, userToken.User.ID, gameID)
}

// GetGameHub retourne le hub d'une partie spécifique
func GetGameHub(gameID string) *websocket.Hub {
	return gameHandler.GetOrCreateGameHub(gameID)
}

// CleanupGameHub supprime le hub d'une partie terminée
func (wsh *WebSocketHandler) CleanupGameHub(gameID string) {
	wsh.mutex.Lock()
	defer wsh.mutex.Unlock()

	delete(wsh.gameHubs, gameID)
}
