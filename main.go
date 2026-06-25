package main

import (
	"log"

	"github.com/moroz/pindakaas/sshserver"
	"golang.org/x/sync/errgroup"
)

func main() {
	var g errgroup.Group

	g.Go(func() error {
		server, err := sshserver.New(2137)
		if err != nil {
			return err
		}

		return server.Serve()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
