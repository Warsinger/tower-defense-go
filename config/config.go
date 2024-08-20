package config

import "github.com/yohamta/donburi"

type ConfigData struct {
	Debug          bool
	GridLines      bool
	ServerPort     string
	ClientHostPort string
	Computer       bool
}

var Config = donburi.NewComponentType[ConfigData]()

func NewConfig(world donburi.World, debug bool, computer bool) *ConfigData {
	entity := world.Create(Config)
	entry := world.Entry(entity)

	Config.Set(entry, &ConfigData{Debug: debug, Computer: computer})
	return Config.Get(entry)
}

func GetConfig(world donburi.World) *ConfigData {
	return Config.Get(Config.MustFirst(world))
}
