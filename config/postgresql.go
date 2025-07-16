package config

import (
	"context"
	"os"

	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v4"
)

func InitPgSQL() (context.Context, *pgx.ConnConfig) {
	ctx := context.Background()
	connstring := "postgresql://"

	if env := os.Getenv("POSTGRES_USER"); env == "" {
		log.Fatal("Bad 'POSTGRES_USER' parameter env")
	} else {
		connstring += env
	}

	if env := os.Getenv("POSTGRES_PASSWORD"); env == "" {
		log.Warn("'POSTGRES_PASSWORD' not set parameter env")
	} else {
		connstring += ":" + env

	}

	if env := os.Getenv("POSTGRES_HOST"); env == "" {
		log.Fatal("Bad 'POSTGRES_HOST' parameter env")
		os.Exit(1)
	} else {
		connstring += "@" + env
	}

	if env := os.Getenv("POSTGRES_DB"); env == "" {
		log.Fatal("Bad 'POSTGRES_DB' parameter env")
	} else {
		connstring += "/" + env
	}

	connstring += "?sslmode=disable"

	connConf, err := pgx.ParseConfig(connstring)
	if err != nil {
		log.Fatalf("Parse error : %s", err)
	}

	sqlCo, err := pgx.ConnectConfig(ctx, connConf)
	if err != nil {
		log.Fatalf("error connect psql : %s", err)
		return ctx, connConf
	}
	defer sqlCo.Close(ctx)

	query := `
	CREATE EXTENSION IF NOT EXISTS pgcrypto;

	CREATE TABLE IF NOT EXISTS account (
		id 								SERIAL,
		email 						TEXT NOT NULL UNIQUE,
		username 					TEXT NOT NULL UNIQUE,
		password 					TEXT NOT NULL,
		recover_token 		TEXT,
		admin 						boolean DEFAULT FALSE,
		enable						boolean DEFAULT TRUE,
		PRIMARY KEY(id)
	);

	-- Table pour les défis entre joueurs
	CREATE TABLE IF NOT EXISTS challenges (
		id 							VARCHAR(36) PRIMARY KEY,
		challenger_id 	INTEGER REFERENCES account(id) NOT NULL,
		challenged_id 	INTEGER REFERENCES account(id) NOT NULL,
		status 					TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'accepted', 'declined', 'expired', 'cancelled')),
		message 				TEXT DEFAULT '',
		game_id 				VARCHAR(36),
		created_at 			TIMESTAMP DEFAULT NOW(),
		updated_at 			TIMESTAMP DEFAULT NOW(),
		expires_at 			TIMESTAMP DEFAULT (NOW() + INTERVAL '24 hours'),
		responded_at 		TIMESTAMP
	);

	-- Table pour les parties (sans référence aux rooms)
	CREATE TABLE IF NOT EXISTS games (
		id 							VARCHAR(36) PRIMARY KEY,
		player1_id 			INTEGER REFERENCES account(id) NOT NULL,
		player2_id 			INTEGER REFERENCES account(id) NOT NULL,
		current_turn 		BIGINT NOT NULL,
		game_phase 			INTEGER DEFAULT 0 CHECK (game_phase IN (0, 1)),
		board 					JSONB DEFAULT '[[null,null,null,null],[null,null,null,null],[null,null,null,null],[null,null,null,null]]',
		available_pieces JSONB DEFAULT '[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15]',
		selected_piece 	INTEGER DEFAULT -1,
		status 					INTEGER DEFAULT 0 CHECK (status IN (0, 1)),
		winner 					BIGINT DEFAULT 0,
		move_history 		JSONB DEFAULT '[]',
		created_at 			TIMESTAMP DEFAULT NOW(),
		updated_at 			TIMESTAMP DEFAULT NOW()
	);

	-- Index pour optimiser les requêtes
	CREATE INDEX IF NOT EXISTS idx_challenges_challenger ON challenges(challenger_id);
	CREATE INDEX IF NOT EXISTS idx_challenges_challenged ON challenges(challenged_id);
	CREATE INDEX IF NOT EXISTS idx_challenges_status ON challenges(status);
	CREATE INDEX IF NOT EXISTS idx_challenges_expires ON challenges(expires_at);
	CREATE INDEX IF NOT EXISTS idx_games_status ON games(status);
	CREATE INDEX IF NOT EXISTS idx_games_player1 ON games(player1_id);
	CREATE INDEX IF NOT EXISTS idx_games_player2 ON games(player2_id);
	CREATE INDEX IF NOT EXISTS idx_games_players ON games(player1_id, player2_id);
	`

	_, err = sqlCo.Exec(ctx, query)
	if err != nil {
		log.Fatal("During postgresql setup", "query", query, "error", err)
	}
	return ctx, connConf
}
