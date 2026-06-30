package sshserver

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/services"
	"github.com/moroz/pindakaas/types"
	"golang.org/x/crypto/ssh"
)

type SSHServer struct {
	serverConfig *ssh.ServerConfig
	hostService  *services.TunnelService
	connRegistry types.TunnelRegistry
}

func New(db queries.DBTX, connRegistry types.TunnelRegistry) (*SSHServer, error) {
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
		hostService:  services.NewTunnelService(db, connRegistry),
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

	// One tunnel per connection, shared between the request-forwarding handler
	// (which fills in the bind details and registers it) and the session
	// handler (which streams its forwarding logs to an interactive client).
	tunnel := &types.Tunnel{Conn: conn}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for req := range reqs {
			if req.Type != "tcpip-forward" {
				if req.WantReply {
					req.Reply(false, nil)
				}
				continue
			}

			var request types.RequestPortForwardingPayload

			err := ssh.Unmarshal(req.Payload, &request)
			if err != nil {
				log.Print("Failed to parse SSH wire format: ", err)
			}

			log.Printf("Forwarding request: %v", request)

			tunnel.BindAddr = request.BindAddr
			tunnel.BindPort = request.BindPort
			tunnel.AllocatedPort = 42069
			s.connRegistry.RegisterConnection(host.Subdomain, tunnel)
			defer s.connRegistry.DeregisterConnection(host.Subdomain, tunnel)

			response := ssh.Marshal(types.RequestPortForwardingSuccessPayload{BindPort: tunnel.AllocatedPort})
			req.Reply(true, response)
		}
	}()

	go func() {
		defer wg.Done()

		for newChan := range chans {
			// Accept interactive sessions (`ssh -tt`) to stream forwarding
			// logs; reject everything else.
			if newChan.ChannelType() != "session" {
				newChan.Reject(ssh.UnknownChannelType, "")
				continue
			}

			go handleSession(newChan, host.Subdomain, tunnel)
		}
	}()

	wg.Wait()

	log.Printf("Connection closed")
}

// handleSession accepts an interactive session channel and streams the tunnel's
// forwarding logs to the client's terminal until the session is closed. It is
// read-only: anything the user types is discarded.
func handleSession(newChan ssh.NewChannel, subdomain string, tunnel *types.Tunnel) {
	channel, requests, err := newChan.Accept()
	if err != nil {
		log.Print("Failed to accept session channel: ", err)
		return
	}
	defer channel.Close()

	// Watch client input for Ctrl-C (0x03) / Ctrl-D (0x04) so the user can
	// disconnect, and otherwise discard keystrokes (this is a read-only log
	// view). There is no PTY on this side to turn ^C into a signal, so we have
	// to close the session ourselves when we see one.
	go func() {
		buf := make([]byte, 256)
		for {
			n, err := channel.Read(buf)
			if err != nil {
				return
			}
			if bytes.ContainsAny(buf[:n], "\x03\x04") {
				// Report that the "command" exited so the client disconnects
				// cleanly on the first keypress. Without an exit-status OpenSSH
				// hangs and needs a second Ctrl-C to abort locally.
				channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{0}))
				channel.Close()
				return
			}
		}
	}()

	lines := make(chan string, 64)
	done := make(chan struct{})

	// Write buffered log lines to the terminal until the session ends. Lines
	// use CRLF because `ssh -tt` puts the client terminal in raw mode.
	go func() {
		for {
			select {
			case line := <-lines:
				fmt.Fprint(channel, line+"\r\n")
			case <-done:
				return
			}
		}
	}()

	attached := false
	for req := range requests {
		switch req.Type {
		case "pty-req", "shell", "exec":
			if req.WantReply {
				req.Reply(true, nil)
			}
			if !attached {
				attached = true
				url := fmt.Sprintf("https://%s.%s", subdomain, config.BaseDomain)
				if config.HTTPSPort != 443 {
					url = net.JoinHostPort(url, strconv.Itoa(int(config.HTTPSPort)))
				}
				fmt.Fprintf(channel, "Streaming forwarding logs for %q. Disconnect with ~. or Ctrl-C.\r\n", url)
				tunnel.AttachSession(channel, lines)
			}
		default:
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}

	tunnel.DetachSession(lines)
	close(done)
}

func (s *SSHServer) authenticateConnection(conn ssh.ConnMetadata) (*ssh.Permissions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	host, err := s.hostService.AuthenticateHostBySSHUsername(ctx, conn.User())
	if err != nil {
		// Returning a BannerError sends the message to the client as an SSH
		// userauth banner (shown by OpenSSH), so the user sees why the
		// connection was rejected instead of a bare "Permission denied".
		return nil, &ssh.BannerError{
			Err:     err,
			Message: "Authentication failed: invalid credentials.\n",
		}
	}

	return &ssh.Permissions{
		ExtraData: map[any]any{
			"host": host,
		},
	}, nil
}
