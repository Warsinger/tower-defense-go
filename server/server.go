package server

import (
	"fmt"
	"log"
	"time"

	"github.com/yohamta/donburi"
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
		return err
	}
	err = s.SendWorld()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) SerializeWorld() error {
	fmt.Println("serializing world...")
	return nil
}
func (s *Server) SendWorld() error {
	fmt.Println("sending world...")
	return nil
}
