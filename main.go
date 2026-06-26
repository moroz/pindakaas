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

	server := httpserver.New(ctx, config.BaseDomain)

	g.Go(func() error {
		return server.ListenAndServe(config.HTTPPort)
	})

	g.Go(func() error {
		return server.ListenAndServeTLS(config.HTTPSPort, config.TLSCertFile, config.TLSKeyFile)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
