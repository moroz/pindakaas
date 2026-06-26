package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/httpserver"
	"github.com/moroz/pindakaas/registry"
	"github.com/moroz/pindakaas/sshserver"
	"golang.org/x/sync/errgroup"

	_ "modernc.org/sqlite"
)

func main() {
	db, _ := sql.Open("sqlite", config.DatabaseUrl)
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer db.Close()

	reg := registry.New()

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		server, err := sshserver.New(db, reg)
		if err != nil {
			return err
		}

		return server.Serve(ctx, config.SSHPort)
	})

	server := httpserver.New(reg)

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
