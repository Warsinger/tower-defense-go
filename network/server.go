package network

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/leap-fish/necs/router"
	"github.com/leap-fish/necs/transports"
	"github.com/yohamta/donburi"
	"nhooyr.io/websocket"
)

const (
	TickRate = 16
)

type Server struct {
	host  transports.NetworkTransport
	world donburi.World
}

func NewServer(world donburi.World, address, port string) *Server {
	portNum, _ := strconv.Atoi(port)
	if address == "" {
		address = "localhost"
	}
	return &Server{
		world: world,
		host: transports.NewWsServerTransport(
			uint(portNum),
			address,
			&websocket.AcceptOptions{
				InsecureSkipVerify: true,
			},
		),
	}
}

func (s *Server) Start() error {
	router.OnConnect(func(sender *router.NetworkClient) {
		fmt.Printf("Client %s connected to the server!\n", sender.Id())
	})
	router.OnDisconnect(func(sender *router.NetworkClient, err error) {
		fmt.Printf("Client %s disconnected from the server! / Reason [%s]\n", sender.Id(), err)
	})

	RegisterComponenets()
	srvsync.UseEsync(s.world)

	go s.StartHost()
	go s.startTicking()

	return nil
}

func (s *Server) StartHost() {
	err := s.host.Start()
	if err != nil {
		log.Fatalf("Error starting host server: %v", err)
	}
}
func (s *Server) startTicking() {
	for range time.NewTicker(time.Second / TickRate).C {
		// fmt.Printf("Syncing world %v...\n", s.world)
		err := srvsync.DoSync()
		if err != nil {
			// TODO don't panic here, just retry later, maybe client will reconnect
			log.Fatalf("Unable to perform esync.DoSync: %v", err)
		}
	}
}
