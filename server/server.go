package server

import (
	"log"
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

func NewServer(world donburi.World, port uint, address string) *Server {
	return &Server{
		world: world,
		host: transports.NewWsServerTransport(
			port,
			address,
			&websocket.AcceptOptions{
				InsecureSkipVerify: true,
			},
		),
	}
}

func (s *Server) Start() {
	srvsync.UseEsync(s.world)

	go s.startTicking()

	err := s.host.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) startTicking() {
	for range time.NewTicker(time.Second / TickRate).C {

		err := srvsync.DoSync()
		if err != nil {
			log.Fatal(err)
		}
	}
}
