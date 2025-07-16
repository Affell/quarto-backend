package challenge

import (
	"fmt"
	"quarto/models/game"
	"time"

	"github.com/google/uuid"
)

// SendChallenge envoie un défi à un autre joueur
func SendChallenge(challengerID, challengedID int64, message string) (*Challenge, error) {
	// Vérifier que le joueur ne se défie pas lui-même
	if challengerID == challengedID {
		return nil, fmt.Errorf("vous ne pouvez pas vous défier vous-même")
	}

	// Vérifier qu'il n'y a pas déjà un défi en attente entre ces joueurs
	existingChallenge, err := GetPendingChallengeBetween(challengerID, challengedID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la vérification des défis existants: %v", err)
	}
	if existingChallenge != nil {
		return nil, fmt.Errorf("un défi est déjà en attente entre ces joueurs")
	}

	// Créer le défi
	challenge := Challenge{
		ID:           uuid.New().String(),
		ChallengerID: challengerID,
		ChallengedID: challengedID,
		Status:       "pending",
		Message:      message,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour), // Expire dans 24h
	}

	// Sauvegarder en base
	err = CreateChallenge(challenge)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création du défi: %v", err)
	}

	return &challenge, nil
}

// AcceptChallenge accepte un défi et crée une partie
func AcceptChallenge(challengeID string, challengedID int64) (*Challenge, *game.Game, error) {
	// Récupérer le défi
	challenge, err := GetChallengeByID(challengeID)
	if err != nil {
		return nil, nil, fmt.Errorf("défi non trouvé: %v", err)
	}

	// Vérifier que c'est le bon joueur qui répond
	if challenge.ChallengedID != challengedID {
		return nil, nil, fmt.Errorf("vous n'êtes pas autorisé à répondre à ce défi")
	}

	// Vérifier que le défi peut être accepté
	if !challenge.CanRespond() {
		return nil, nil, fmt.Errorf("ce défi ne peut plus être accepté (expiré ou déjà traité)")
	}

	// Créer une nouvelle partie
	newGame, err := game.CreateNewGame(challenge.ChallengerID, challenge.ChallengedID)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur lors de la création de la partie: %v", err)
	}

	// Mettre à jour le défi
	err = UpdateChallengeStatus(challengeID, "accepted", &newGame.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur lors de l'acceptation du défi: %v", err)
	}

	// Récupérer le défi mis à jour
	updatedChallenge, err := GetChallengeByID(challengeID)
	if err != nil {
		return nil, nil, fmt.Errorf("erreur lors de la récupération du défi mis à jour: %v", err)
	}

	return updatedChallenge, &newGame, nil
}

// DeclineChallenge refuse un défi
func DeclineChallenge(challengeID string, challengedID int64) (*Challenge, error) {
	// Récupérer le défi
	challenge, err := GetChallengeByID(challengeID)
	if err != nil {
		return nil, fmt.Errorf("défi non trouvé: %v", err)
	}

	// Vérifier que c'est le bon joueur qui répond
	if challenge.ChallengedID != challengedID {
		return nil, fmt.Errorf("vous n'êtes pas autorisé à répondre à ce défi")
	}

	// Vérifier que le défi peut être refusé
	if !challenge.CanRespond() {
		return nil, fmt.Errorf("ce défi ne peut plus être refusé (expiré ou déjà traité)")
	}

	// Mettre à jour le défi
	err = UpdateChallengeStatus(challengeID, "declined", nil)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du refus du défi: %v", err)
	}

	// Récupérer le défi mis à jour
	updatedChallenge, err := GetChallengeByID(challengeID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du défi mis à jour: %v", err)
	}

	return updatedChallenge, nil
}

// CancelChallenge annule un défi envoyé
func CancelChallenge(challengeID string, challengerID int64) (*Challenge, error) {
	// Récupérer le défi
	challenge, err := GetChallengeByID(challengeID)
	if err != nil {
		return nil, fmt.Errorf("défi non trouvé: %v", err)
	}

	// Vérifier que c'est le bon joueur qui annule
	if challenge.ChallengerID != challengerID {
		return nil, fmt.Errorf("vous n'êtes pas autorisé à annuler ce défi")
	}

	// Vérifier que le défi peut être annulé
	if challenge.Status != "pending" {
		return nil, fmt.Errorf("ce défi ne peut plus être annulé")
	}

	// Mettre à jour le défi
	err = UpdateChallengeStatus(challengeID, "cancelled", nil)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'annulation du défi: %v", err)
	}

	// Récupérer le défi mis à jour
	updatedChallenge, err := GetChallengeByID(challengeID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du défi mis à jour: %v", err)
	}

	return updatedChallenge, nil
}

// GetMyChallenges récupère tous les défis d'un utilisateur organisés par type
func GetMyChallenges(userID int64) (*ChallengeListResponse, error) {
	challenges, err := GetUserChallenges(userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des défis: %v", err)
	}

	response := &ChallengeListResponse{
		Sent:     make([]Challenge, 0),
		Received: make([]Challenge, 0),
	}

	for _, challenge := range challenges {
		if challenge.ChallengerID == userID {
			response.Sent = append(response.Sent, challenge)
		} else {
			response.Received = append(response.Received, challenge)
		}
	}

	return response, nil
}
