package config

import "github.com/yohamta/donburi"

type ConfigData struct {
	Computer  bool
	Debug     bool
	GridLines bool
	ShowStats bool
	Sound     bool

	ClientHostPort string
	ServerPort     string
}

var Config = donburi.NewComponentType[ConfigData]()

func NewConfig(world donburi.World, debug, computer, sound bool) *ConfigData {
	entity := world.Create(Config)
	entry := world.Entry(entity)

	Config.Set(entry, &ConfigData{Debug: debug, Computer: computer, Sound: sound})
	return Config.Get(entry)
}

func GetConfig(world donburi.World) *ConfigData {
	return Config.Get(Config.MustFirst(world))
}
