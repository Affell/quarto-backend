package main

import (
	"quarto/config"
	"quarto/models/postgresql"
	"time"

	"github.com/charmbracelet/log"
)

func init() {
	start := time.Now()

	config.Init(Folder)
	postgresql.SQLCtx, postgresql.SQLConn = config.InitPgSQL()

	log.Debug("Initialization ended", "took", time.Since(start).Round(time.Millisecond).String())
}
