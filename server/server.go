package server

import (
	"fmt"
	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
	"log"
	"net"
	"time"
)

const (
	TickRate = 16
)

type Server struct {
	host  string
	world donburi.World
}

func NewServer(world donburi.World, port uint, address string) *Server {
	return &Server{
		world: world,
		host:  fmt.Sprintf("%v:%v", address, port),
	}
}

func (s *Server) Start() {
	srvsync.UseEsync(s.world)

	go s.startTicking()

	// TODO start listening
}

func (s *Server) startTicking() {
	for range time.NewTicker(time.Second / TickRate).C {

		err := s.DoSync()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (s *Server) DoSync() error {
	fmt.Println("syncing...")
	err := s.SerializeWorld()
	if err != nil {
		log.Fatal(err)
	}
	err = s.SendWorld()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) SerializeWorld() error {
	fmt.Println("serializing world...")
	return nil
}
func (s *Server) SendWorld() error {
	fmt.Println("sending world...")
	return nil
}
