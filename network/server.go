package network

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/yohamta/donburi"
)

const (
	TickRate = 16
)

type Server struct {
	host     string
	world    donburi.World
	listener net.Listener
	conn     net.Conn
}

func NewServer(world donburi.World, host, port string) *Server {
	return &Server{
		world: world,
		host:  net.JoinHostPort(host, port),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.host)
	if err != nil {
		return err
	}
	s.listener = listener
	fmt.Printf("listening %v\n", s.host)

	go func() error {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		s.conn = conn

		go func(conn net.Conn) {
			data, err := io.ReadAll(conn)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("data after connection %v\n", string(data))
			conn.Write([]byte("Connection Success\n"))
		}(conn)
		return nil
	}()
	// go s.startTicking()
	return nil
}

func (s *Server) Close() {
	s.conn.Close()
	s.listener.Close()
}

// func (s *Server) startTicking() {
// 	for range time.NewTicker(time.Second / TickRate).C {

// 		err := s.DoSync()
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

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
