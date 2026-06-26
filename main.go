package main

import (
	"context"
	"log"

	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/httpserver"
	"github.com/moroz/pindakaas/sshserver"
	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		server, err := sshserver.New(ctx, config.SSHPort)
		if err != nil {
			return err
		}

		return server.Serve()
	})

	g.Go(func() error {
		server := httpserver.New(ctx, config.HTTPPort)
		return server.Serve()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
