package pokemon

import (
	"github.com/olebedev/config"
	"github.com/doctorfree/wtf/cfg"
)

const (
	defaultFocusable = false
	defaultTitle     = "Pokémon"
)

type Settings struct {
	*cfg.Common

	id       int
	random   bool
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		id:       ymlConfig.UInt("id", 0),
		random:   ymlConfig.UBool("random", true),
	}

	settings.SetDocumentationPath("pokemon")

	return &settings
}
