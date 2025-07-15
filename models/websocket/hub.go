package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // En production, vérifier l'origine
	},
}

type Hub struct {
	clients     map[*Client]bool            // Tous les clients connectés
	gameClients map[string]map[*Client]bool // gameID -> clients de cette partie
	register    chan *Client
	unregister  chan *Client
	mutex       sync.RWMutex
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID int64
	gameID string
}

type WSMessage struct {
	Type   string `json:"type"`
	GameID string `json:"game_id,omitempty"`
	UserID string `json:"user_id"`
	Data   any    `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		gameClients: make(map[string]map[*Client]bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true

	// Ajouter à la partie
	if client.gameID != "" {
		if h.gameClients[client.gameID] == nil {
			h.gameClients[client.gameID] = make(map[*Client]bool)
		}
		h.gameClients[client.gameID][client] = true
	}

	log.Printf("Client connecté: %d dans la partie %s", client.userID, client.gameID)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// Retirer de la partie
		if client.gameID != "" && h.gameClients[client.gameID] != nil {
			if _, ok := h.gameClients[client.gameID][client]; ok {
				delete(h.gameClients[client.gameID], client)

				// Supprimer la partie si vide
				if len(h.gameClients[client.gameID]) == 0 {
					delete(h.gameClients, client.gameID)
				}
			}
		}

		log.Printf("Client déconnecté: %d", client.userID)
	}
}

// BroadcastToGame envoie un message à tous les joueurs d'une partie
func (h *Hub) BroadcastToGame(gameID string, message WSMessage) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	log.Printf("BroadcastToGame called for game: %s, message type: %s", gameID, message.Type)

	if gameClients, exists := h.gameClients[gameID]; exists {
		log.Printf("Found %d clients in game %s", len(gameClients), gameID)

		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("Erreur de sérialisation du message: %v", err)
			return
		}

		sentCount := 0
		for client := range gameClients {
			select {
			case client.send <- messageBytes:
				sentCount++
				log.Printf("Message sent to client %d", client.userID)
			default:
				log.Printf("Failed to send message to client %d, removing from game", client.userID)
				close(client.send)
				delete(h.clients, client)
				delete(gameClients, client)
			}
		}
		log.Printf("Successfully sent message to %d/%d clients in game %s", sentCount, len(gameClients), gameID)
	} else {
		log.Printf("Game %s not found or has no clients", gameID)
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var message WSMessage
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Erreur WebSocket: %v", err)
			}
			break
		}

		// Traiter le message reçu
		c.handleMessage(message)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(message WSMessage) {
	log.Printf("Message reçu de %d: %s", c.userID, message.Type)

	switch message.Type {
	case "ping":
		response := WSMessage{
			Type:   "pong",
			UserID: "server",
			Data:   map[string]string{"message": "pong"},
		}
		responseBytes, _ := json.Marshal(response)
		c.send <- responseBytes

	default:
		log.Printf("Type de message non géré: %s", message.Type)
	}
}

// HandleWebSocket gère la connexion WebSocket
func (h *Hub) HandleWebSocket(c echo.Context, userID int64, gameID string) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Erreur de mise à niveau WebSocket: %v", err)
		return err
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
		gameID: gameID,
	}

	client.hub.register <- client

	// Démarrer les goroutines pour lire et écrire
	go client.writePump()
	go client.readPump()

	return nil
}
