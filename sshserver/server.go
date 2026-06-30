package sshserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/registry"
	"github.com/moroz/pindakaas/services"
	"github.com/moroz/pindakaas/types"
	"golang.org/x/crypto/ssh"
)

type SSHServer struct {
	serverConfig *ssh.ServerConfig
	hostService  *services.HostService
	connRegistry *registry.Registry
}

func New(db queries.DBTX, connRegistry *registry.Registry) (*SSHServer, error) {
	algorithms := ssh.SupportedAlgorithms()

	serverConfig := &ssh.ServerConfig{
		Config: ssh.Config{
			KeyExchanges: algorithms.KeyExchanges,
			MACs:         algorithms.MACs,
			Ciphers:      algorithms.Ciphers,
		},
		NoClientAuth: true,
	}

	privateBytes, err := os.ReadFile(config.SSHServerKeyPath)
	if err != nil {
		return nil, fmt.Errorf("SSHServer.Serve: Failed to load private key: %w", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("SSHServer.Serve: Failed to parse private key: %w", err)
	}

	serverConfig.AddHostKey(private)

	server := &SSHServer{
		serverConfig: serverConfig,
		hostService:  services.NewHostService(db),
		connRegistry: connRegistry,
	}

	serverConfig.NoClientAuthCallback = server.authenticateConnection

	return server, nil
}

func (s *SSHServer) Serve(ctx context.Context, port uint16) error {
	listenOn := config.FormatHostPort(port)
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("Failed to bind on port %v: %w", port, err)
	}

	log.Printf("SSH server listening on %s", listener.Addr())

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Print("Failed to accept incoming connection: ", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *SSHServer) handleConn(newConnection net.Conn) {
	conn, chans, reqs, err := ssh.NewServerConn(newConnection, s.serverConfig)
	if err != nil {
		log.Print("SSH handshake failed: ", err)
		return
	}
	defer conn.Close()

	host := conn.Permissions.ExtraData["host"].(*queries.Tunnel)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for req := range reqs {
			if req.Type != "tcpip-forward" && req.WantReply {
				log.Printf("Rejecting request %s", req.Type)
				req.Reply(false, nil)
				continue
			}

			var request types.RequestPortForwardingPayload

			err := ssh.Unmarshal(req.Payload, &request)
			if err != nil {
				log.Print("Failed to parse SSH wire format: ", err)
			}

			log.Printf("Forwarding request: %v", request)

			tunnel := &types.Tunnel{
				Conn:          conn,
				BindAddr:      request.BindAddr,
				BindPort:      request.BindPort,
				AllocatedPort: 42069,
			}
			s.connRegistry.RegisterConnection(host.Subdomain, tunnel)
			defer s.connRegistry.DeregisterConnection(host.Subdomain, tunnel)

			response := ssh.Marshal(types.RequestPortForwardingSuccessPayload{BindPort: tunnel.AllocatedPort})
			req.Reply(true, response)
		}

		wg.Done()
	}()

	// Reject all channels for now
	go func() {
		for newChan := range chans {
			newChan.Reject(ssh.UnknownChannelType, "")
		}

		wg.Done()
	}()

	wg.Wait()

	log.Printf("Connection closed")
}

func (s *SSHServer) authenticateConnection(conn ssh.ConnMetadata) (*ssh.Permissions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	host, err := s.hostService.AuthenticateHostBySSHUsername(ctx, conn.User())
	if err != nil {
		return nil, err
	}

	return &ssh.Permissions{
		ExtraData: map[any]any{
			"host": host,
		},
	}, nil
}
