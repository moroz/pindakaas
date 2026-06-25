package sshserver

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"

	"github.com/moroz/pindakaas/config"
	"golang.org/x/crypto/ssh"
)

type SSHServer struct {
	Port         uint16
	ServerConfig *ssh.ServerConfig
}

func New(port uint16) (*SSHServer, error) {
	algorithms := ssh.SupportedAlgorithms()

	serverConfig := &ssh.ServerConfig{
		Config: ssh.Config{
			KeyExchanges: algorithms.KeyExchanges,
			MACs:         algorithms.MACs,
			Ciphers:      algorithms.Ciphers,
		},
		NoClientAuth: true,
		NoClientAuthCallback: func(conn ssh.ConnMetadata) (*ssh.Permissions, error) {
			log.Printf("%+v", conn)
			return &ssh.Permissions{}, nil
		},
	}

	privateBytes, err := os.ReadFile(config.ServerKeyPath)
	if err != nil {
		return nil, fmt.Errorf("SSHServer.Serve: Failed to load private key: %w", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("SSHServer.Serve: Failed to parse private key: %w", err)
	}

	serverConfig.AddHostKey(private)

	return &SSHServer{
		Port:         port,
		ServerConfig: serverConfig,
	}, nil
}

func (s *SSHServer) Serve() error {

	listener, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(int(s.Port))))
	if err != nil {
		return fmt.Errorf("Failed to bind on port %v: %w", s.Port, err)
	}

	log.Printf("SSH server listening on %s", listener.Addr())

	for {
		go s.handleConn(listener.Accept())
	}
}

func (s *SSHServer) handleConn(newConnection net.Conn, err error) {
	if err != nil {
		log.Print("Failed to accept incoming connection: ", err)
		return
	}

	conn, chans, reqs, err := ssh.NewServerConn(newConnection, s.ServerConfig)
	if err != nil {
		log.Print("SSH handshake failed: ", err)
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for req := range reqs {
			if req.Type != "tcpip-forward" && req.WantReply {
				log.Printf("Rejecting request %s", req.Type)
				req.Reply(false, nil)
				continue
			}

			var request struct {
				BindAddr string
				BindPort uint32
			}

			err := ssh.Unmarshal(req.Payload, &request)
			if err != nil {
				log.Print("Failed to parse SSH wire format: ", err)
			}

			log.Printf("Forwarding request: %v", request)

			response := ssh.Marshal(struct{ Port uint32 }{0})
			req.Reply(true, response)
		}

		wg.Done()
	}()

	go func() {
		for newChan := range chans {
			log.Printf("Channel type: %s", newChan.ChannelType())
			newChan.Reject(ssh.UnknownChannelType, "")
		}

		wg.Done()
	}()

	wg.Wait()

	log.Printf("Connection closed")
}
