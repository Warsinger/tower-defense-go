package config

import "github.com/yohamta/donburi"

type ConfigData struct {
	debug     bool
	gridLines bool
}

var Config = donburi.NewComponentType[ConfigData]()

func NewConfig(w donburi.World, debug bool) *ConfigData {
	entity := w.Create(Config)
	entry := w.Entry(entity)

	Config.SetValue(entry, ConfigData{debug: debug, gridLines: false})
	return Config.Get(entry)
}

func (c *ConfigData) IsDebug() bool {
	return c.debug
}

func (c *ConfigData) SetDebug(debug bool) {
	c.debug = debug
}

func (c *ConfigData) IsGridLines() bool {
	return c.gridLines
}

func (c *ConfigData) SetGridLines(gridLines bool) {
	c.gridLines = gridLines
}
