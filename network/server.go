package network

import (
	"log"
	"strconv"
	"time"

	"github.com/leap-fish/necs/esync/srvsync"
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
	return &Server{
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
	RegisterComponenets()
	srvsync.UseEsync(s.world)

	go s.startTicking()
	go s.StartHost()

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
		err := srvsync.DoSync()
		if err != nil {
			log.Fatalf("Unable to perform esync.DoSync: %v", err)
		}
	}
}
