package network

import (
	"bufio"
	"fmt"
	"net"

	"github.com/yohamta/donburi"
)

type Client struct {
	host  string
	world donburi.World
	conn  net.Conn
}

func NewClient(world donburi.World, host, port string) *Client {
	return &Client{
		world: world,
		host:  net.JoinHostPort(host, port),
	}
}

func (c *Client) Start() error {
	fmt.Printf("Connecting...\n")
	// connect to server
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return err
	}
	c.conn = conn

	fmt.Printf("Connected\n")
	c.SendMessage("Connecting\n")
	return nil
}

func (c *Client) SendMessage(message string) error {
	// BUG something here blocks, is server responding? should excecute in goroutine
	fmt.Fprint(c.conn, message)
	status, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("return status %v\n", status)
	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}
