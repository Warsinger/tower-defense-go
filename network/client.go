package network

import (
	"context"

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
}

func NewClient(world donburi.World, address string) *Client {
	return &Client{
		World:     world,
		Transport: transports.NewWsClientTransport(address),
	}
}
func NewClientNewWorld(address string) *Client {
	return &Client{
		World:     donburi.NewWorld(),
		Transport: transports.NewWsClientTransport(address),
	}
}

func (c *Client) Start() error {
	RegisterComponenets()
	clisync.RegisterClient(c.World)
	c.Transport.Start(func(conn *websocket.Conn) {
		c.Network = router.NewNetworkClient(context.Background(), conn)
	})
	return nil
}
