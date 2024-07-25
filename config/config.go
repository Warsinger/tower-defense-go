package config

import "github.com/yohamta/donburi"

type ConfigData struct {
	debug bool
}

var Config = donburi.NewComponentType[ConfigData]()

func NewConfig(w donburi.World, debug bool) *ConfigData {
	entity := w.Create(Config)
	entry := w.Entry(entity)

	Config.SetValue(entry, ConfigData{debug: debug})
	return Config.Get(entry)
}

func (c *ConfigData) IsDebug() bool {
	return c.debug
}
func (c *ConfigData) SetDebug(debug bool) {
	c.debug = debug
}
