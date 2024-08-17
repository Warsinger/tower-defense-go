package scenes

import (
	"fmt"
	"strings"
	"tower-defense/config"
	"tower-defense/network"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leap-fish/necs/router"
	"github.com/yohamta/donburi"
)

type Controller struct {
	server *network.Server
	client *network.Client
}

const urlPrefix = "ws://"

func (c *Controller) StartServer(world donburi.World, gameOptions *config.ConfigData, newGameCallback NewGameCallback) error {
	ebiten.SetWindowTitle("Tower Defense (server)")
	fmt.Printf("listening on port %v\n", gameOptions.ServerPort)
	c.server = network.NewServer(world, "", gameOptions.ServerPort)
	err := c.server.Start()
	if err != nil {
		return err
	}
	if c.client != nil {
		// TODO stop client
		c.client = nil
	}
	router.On(func(sender *router.NetworkClient, message network.ClientConnectMessage) {
		fmt.Println("recv client connect message")
		c.startClient(nil, message.Address)
	})
	registerStartGame(newGameCallback, gameOptions)

	return nil
}

func (c *Controller) StartClient(world donburi.World, gameOptions *config.ConfigData, newGameCallback NewGameCallback) error {
	ebiten.SetWindowTitle("Tower Defense (client)")
	clientHostPort := gameOptions.ClientHostPort
	if !strings.HasPrefix(clientHostPort, urlPrefix) {
		clientHostPort = urlPrefix + clientHostPort
	}
	err := c.startClient(world, clientHostPort)
	if err != nil {
		return err
	}
	if c.server != nil {
		// TODO stop server
		c.server = nil
	}
	registerStartGame(newGameCallback, gameOptions)
	return nil
}

func registerStartGame(newGameCallback NewGameCallback, gameOptions *config.ConfigData) {
	router.On(func(sender *router.NetworkClient, message network.StartGameMessage) {
		fmt.Println("recv start game message")
		newGameCallback(false, &controller, gameOptions)
	})

}

func (c *Controller) startClient(world donburi.World, address string) error {
	fmt.Printf("connect to %v\n", address)
	var err error
	// TODO stop existing client
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

func (c *Controller) GetClientWorld() donburi.World {
	if c.client != nil && c.client.World != nil {
		return c.client.World
	}
	return nil
}
