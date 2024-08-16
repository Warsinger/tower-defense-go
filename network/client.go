package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"

	"github.com/leap-fish/necs/esync/clisync"
	"github.com/leap-fish/necs/router"
	"github.com/leap-fish/necs/transports"
	"github.com/yohamta/donburi"
	"nhooyr.io/websocket"
)

type Client struct {
	Transport *transports.WsClientTransport
	Network   *router.NetworkClient
	World     donburi.World
	port      string
}

func NewClientNewWorld(address string) (*Client, error) {

	url, err := url.ParseRequestURI(address)
	if err != nil {
		fmt.Printf("Error parsing address: %s, %v\n", address, err)
		return nil, err
	}
	port := url.Port()
	return &Client{
		World:     donburi.NewWorld(),
		Transport: transports.NewWsClientTransport(address),
		port:      port,
	}, nil
}

func (c *Client) Start(world donburi.World) error {
	if world != nil {
		router.OnConnect(func(client *router.NetworkClient) {
			fmt.Println("Client connected, starting server")
			// start server on a port one higher than the port we connected on and send the address to the server so it can connect back to us
			port, err := strconv.Atoi(c.port)
			if err != nil {
				log.Fatalf("Error getting port: %v", err)
			}
			port++
			address, err := getIPAddress()
			if err != nil {
				log.Fatalf("Error getting address: %v", err)
			}
			server := NewServer(world, address, c.port)
			err = server.Start()
			if err != nil {
				log.Fatalf("Error starting server: %v", err)
			}
			connStr := fmt.Sprintf("ws://%s:%d", address, port)
			c.Network.SendMessage(ClientConnectMessage{connStr})
		})
	}
	RegisterComponenets()
	clisync.RegisterClient(c.World)
	fmt.Println("registered world")
	go func() {
		c.Transport.Start(func(conn *websocket.Conn) {
			fmt.Println("starting client connection")
			c.Network = router.NewNetworkClient(context.Background(), conn)
		})
	}()

	return nil
}

func getIPAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var ip net.IP
	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			if !v.IP.IsLoopback() {
				if v.IP.To4() != nil { //Verify if IP is IPV4
					ip = v.IP
				}
			}
		}
	}
	if ip != nil {
		return ip.String(), nil
	} else {
		return "", err
	}
}
