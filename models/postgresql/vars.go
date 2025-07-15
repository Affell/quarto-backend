package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4"
)

var (
	SQLCtx  context.Context
	SQLConn *pgx.ConnConfig
)
