package pokemon

import (
	"github.com/olebedev/config"
	"github.com/doctorfree/wtf/cfg"
)

const (
	defaultFocusable = false
	defaultTitle     = "Pokémon"
)

type colors struct {
	name  string
	value string
}

type Settings struct {
	*cfg.Common

	colors
	id       int
	random   bool
	language string
	args []interface{} `help:"Defines what data to display and the order." values:"'height', 'weight', 'name', 'genus', and/or 'id'"`
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		id:       ymlConfig.UInt("id", 0),
		random:   ymlConfig.UBool("random", true),
		language: ymlConfig.UString("language", "en"),
		args:     ymlConfig.UList("args"),
	}

	settings.colors.name = ymlConfig.UString("colors.name", "red")
	settings.colors.value = ymlConfig.UString("colors.value", "white")

	settings.SetDocumentationPath("pokemon")

	return &settings
}
