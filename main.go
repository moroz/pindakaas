package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/httpserver"
	"github.com/moroz/pindakaas/registry"
	"github.com/moroz/pindakaas/sshserver"
	"github.com/moroz/pindakaas/web/handlers"
	"github.com/moroz/pindakaas/web/sessions"
	"golang.org/x/sync/errgroup"

	_ "modernc.org/sqlite"
)

func main() {
	db, _ := sql.Open("sqlite", config.DatabaseUrl)
	if _, err := db.ExecContext(context.Background(), "pragma foreign_keys = on"); err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	reg := registry.New()
	sessionStore, err := sessions.NewStore(config.SessionKey)
	if err != nil {
		log.Fatal(err)
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		server, err := sshserver.New(db, reg)
		if err != nil {
			return err
		}

		return server.Serve(ctx, config.SSHPort)
	})

	server := httpserver.New(&handlers.RouterProps{
		DB:             db,
		Store:          sessionStore,
		TunnelRegistry: reg,
	})

	g.Go(func() error {
		return server.ListenAndServe(ctx, config.HTTPPort)
	})

	g.Go(func() error {
		return server.ListenAndServeTLS(ctx, config.HTTPSPort, config.TLSCertFile, config.TLSKeyFile)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
