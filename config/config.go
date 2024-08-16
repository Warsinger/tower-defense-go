package config

import "github.com/yohamta/donburi"

type ConfigData struct {
	Debug          bool
	GridLines      bool
	ServerPort     string
	ClientHostPort string
}

var Config = donburi.NewComponentType[ConfigData]()

func NewConfig(world donburi.World, debug bool) *ConfigData {
	entity := world.Create(Config)
	entry := world.Entry(entity)

	Config.Set(entry, &ConfigData{Debug: debug})
	return Config.Get(entry)
}

func GetConfig(world donburi.World) *ConfigData {
	return Config.Get(Config.MustFirst(world))
}
