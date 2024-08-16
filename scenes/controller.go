package scenes

import (
	"fmt"
	"strings"
	"tower-defense/network"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leap-fish/necs/router"
	"github.com/yohamta/donburi"
)

type Controller struct {
	server *network.Server
	client *network.Client
}

func (c *Controller) StartServer(world donburi.World, gameOptions *GameOptions, newGameCallback NewGameCallback) error {
	ebiten.SetWindowTitle("Tower Defense (server)")
	fmt.Printf("listening on port %v\n", gameOptions.serverPort)
	c.server = network.NewServer(world, "", gameOptions.serverPort)
	err := c.server.Start()
	if err != nil {
		return err
	}
	if c.client != nil {
		// TODO stop client
	}
	router.On(func(sender *router.NetworkClient, message network.ClientConnectMessage) {
		fmt.Println("recv client connect message")
		c.startClient(nil, message.Address)
	})
	registerStartGame(newGameCallback, gameOptions)

	return nil
}

func (c *Controller) StartClient(world donburi.World, gameOptions *GameOptions, newGameCallback NewGameCallback) error {
	ebiten.SetWindowTitle("Tower Defense (client)")
	clientHostPort := gameOptions.clientHostPort
	if !strings.HasPrefix(clientHostPort, urlPrefix) {
		clientHostPort = urlPrefix + clientHostPort
	}
	err := c.startClient(world, clientHostPort)
	if err != nil {
		return err
	}
	if c.server != nil {
		// TODO stop server
	}
	registerStartGame(newGameCallback, gameOptions)
	return nil
}

func registerStartGame(newGameCallback NewGameCallback, gameOptions *GameOptions) {
	router.On(func(sender *router.NetworkClient, message network.StartGameMessage) {
		fmt.Println("recv start game message")
		newGameCallback(false, gameOptions)
	})

}

func (c *Controller) startClient(world donburi.World, address string) error {
	fmt.Printf("connect to %v\n", address)
	var err error
	c.client, err = network.NewClientNewWorld(address)
	if err != nil {
		return err
	}
	err = c.client.Start(world)
	if err != nil {
		return err
	}
	return nil
}
