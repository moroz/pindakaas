package main

import (
	"log"
	"net"
	"os"

	"github.com/moroz/pindakaas/config"
	"golang.org/x/crypto/ssh"
)

func main() {
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
		log.Fatal("Failed to load private key: ", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	serverConfig.AddHostKey(private)

	listener, err := net.Listen("tcp", "0.0.0.0:2137")
	if err != nil {
		log.Fatal("Failed to bind on port 2137: ", err)
	}

	log.Printf("Listening on port 2137")

	go func() {
		for {
			nConn, err := listener.Accept()
			if err != nil {
				log.Print("Failed to accept incoming connection: ", err)
				continue
			}

			go func() {
				conn, chans, reqs, err := ssh.NewServerConn(nConn, serverConfig)
				if err != nil {
					log.Print("SSH handshake failed: ", err)
				}
				done := make(chan bool, 1)

				defer conn.Close()

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

					done <- true
				}()

				go func() {
					for newChan := range chans {
						log.Printf("Channel type: %s", newChan.ChannelType())
						newChan.Reject(ssh.UnknownChannelType, "")
					}

					done <- true
				}()

				<-done
			}()
		}
	}()

	select {} // block forever
}
