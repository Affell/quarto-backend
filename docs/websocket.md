# WebSocket Documentation - Quarto Backend

## Vue d'ensemble

Le système WebSocket du backend Quarto est conçu exclusivement pour la synchronisation en temps réel des parties de jeu. Il utilise une architecture **hub par partie** simplifiée pour ne gérer que les mises à jour de gameplay.

## Architecture

### Principe

- **Un hub par partie** : Chaque partie active dispose de son propre hub WebSocket
- **Création à la demande** : Les hubs sont créés automatiquement lors de la première connexion à une partie
- **Gameplay uniquement** : Les WebSockets ne servent qu'aux mises à jour de coups en temps réel
- **Nettoyage automatique** : Les hubs vides sont supprimés automatiquement lors de la déconnexion du dernier joueur

### Structure simplifiée

```
Hub
├── clients: map[*Client]bool                    // Tous les clients connectés
├── gameClients: map[gameID]map[*Client]bool     // Clients par partie
├── register/unregister: chan *Client            // Canaux de gestion des connexions
└── mutex: sync.RWMutex                          // Protection concurrentielle
```

## Messages WebSocket

### Format général

```json
{
  "type": "string",
  "game_id": "string",
  "user_id": "string",
  "data": {}
}
```

### Messages entrants (Client → Serveur)

#### ping

Test de connexion et maintien de la session active.

```json
{
  "type": "ping",
  "user_id": "123",
  "data": {}
}
```

### Messages sortants (Serveur → Client)

#### pong

Réponse au ping pour confirmer la connexion.

```json
{
  "type": "pong",
  "user_id": "server",
  "data": {
    "message": "pong"
  }
}
```

#### piece_selected

Une pièce a été sélectionnée par un joueur.

```json
{
  "type": "piece_selected",
  "game_id": "abc-123-def",
  "user_id": "123",
  "data": {
    "id": "abc-123-def",
    "current_turn": "player2",
    "game_phase": "placePiece",
    "selected_piece": 5,
    "board": "[[null,null,null,null]...]",
    "available_pieces": "[0,1,2,3,4,6,7...]"
    // ... état complet de la partie
  }
}
```

#### piece_placed

Une pièce a été placée sur le plateau.

```json
{
  "type": "piece_placed",
  "game_id": "abc-123-def",
  "user_id": "456",
  "data": {
    "id": "abc-123-def",
    "current_turn": "player1",
    "game_phase": "selectPiece",
    "selected_piece": null,
    "board": "[[5,null,null,null]...]",
    "move_history": "[\"BCGP-a1\"]"
    // ... état complet de la partie
  }
}
```

#### game_finished

La partie s'est terminée (victoire ou match nul).

```json
{
  "type": "game_finished",
  "game_id": "abc-123-def",
  "user_id": "123",
  "data": {
    "id": "abc-123-def",
    "status": "finished",
    "winner": "player1"
    // ... état final de la partie
  }
}
```

#### game_forfeited

Un joueur a abandonné la partie.

```json
{
  "type": "game_forfeited",
  "game_id": "abc-123-def",
  "user_id": "456",
  "data": {
    "id": "abc-123-def",
    "status": "finished",
    "winner": "player1"
    // ... état final de la partie
  }
}
```

## Flux d'utilisation

### 1. Connexion à une partie

```javascript
// Les deux joueurs se connectent
const ws1 = new WebSocket(
  "ws://localhost:8080/ws?user_id=123&game_id=game-456"
);
const ws2 = new WebSocket(
  "ws://localhost:8080/ws?user_id=789&game_id=game-456"
);
```

### 2. Écoute des événements

```javascript
ws1.onmessage = (event) => {
  const message = JSON.parse(event.data);

  switch (message.type) {
    case "piece_selected":
      // Mettre à jour l'interface : pièce sélectionnée
      updateSelectedPiece(message.data.selected_piece);
      break;

    case "piece_placed":
      // Mettre à jour le plateau de jeu
      updateBoard(message.data.board);
      updateAvailablePieces(message.data.available_pieces);
      break;

    case "game_finished":
      // Afficher le résultat de la partie
      showGameResult(message.data.winner);
      break;
  }
};
```

### 3. Envoi de ping (optionnel)

```javascript
// Test de connexion périodique
setInterval(() => {
  ws1.send(
    JSON.stringify({
      type: "ping",
      user_id: "123",
      data: {},
    })
  );
}, 30000); // Toutes les 30 secondes
```

## Gestion côté serveur

### Utilisation dans les handlers

```go
// Obtenir le hub d'une partie
hub := websocketHandler.GetGameHub(gameID)

// Diffuser un message à tous les joueurs de la partie
message := websocket.WSMessage{
    Type:   "piece_selected",
    GameID: gameID,
    UserID: strconv.FormatInt(userID, 10),
    Data:   gameData,
}
hub.BroadcastToGame(gameID, message)
```

### Nettoyage automatique

```go
// Le nettoyage des hubs vides se fait automatiquement
// Pas besoin d'appel manuel - le hub se supprime quand le dernier client se déconnecte
```

## Bonnes pratiques

### Côté client

- **Reconnexion automatique** : Implémenter une logique de reconnexion en cas de déconnexion
- **Gestion d'état** : Synchroniser l'état local avec les messages reçus
- **Timeout** : Gérer les timeouts de connexion
- **Fermeture propre** : Fermer la connexion WebSocket à la fin de la partie

### Côté serveur

- **Validation** : Toujours valider les IDs utilisateur et partie
- **Nettoyage** : Nettoyer les hubs des parties terminées pour libérer la mémoire
- **Logs** : Logger les connexions/déconnexions pour le debug

## Sécurité

- **Authentification** : L'`user_id` doit être validé côté serveur
- **Autorisation** : Vérifier que l'utilisateur a le droit d'accéder à la partie
- **Origine** : En production, vérifier l'origine des connexions WebSocket
- **Rate limiting** : Limiter le nombre de messages par client

## Monitoring

### Métriques à surveiller

- Nombre de hubs actifs
- Nombre de connexions par hub
- Fréquence des messages
- Latence des messages
- Erreurs de connexion

### Logs importants

- Connexions/déconnexions WebSocket
- Création/suppression de hubs
- Erreurs de diffusion de messages
- Messages malformés ou non autorisés
