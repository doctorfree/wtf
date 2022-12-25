package pokemon

import (
	"github.com/olebedev/config"
	"github.com/doctorfree/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "Pokémon"
)

type colors struct {
	name  string
	value string
}

type Settings struct {
	*cfg.Common

	colors
	pokemon_id   int
	pokemon_name string
	random       bool
	language     string
	attributes []interface{} `help:"Defines what data to display and the order." values:"'size', 'experience', 'genus', and/or 'text'"`
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		pokemon_id:   ymlConfig.UInt("pokemon_id", 0),
		pokemon_name: ymlConfig.UString("pokemon_name", ""),
		random:       ymlConfig.UBool("random", true),
		language:     ymlConfig.UString("language", "en"),
		attributes:   ymlConfig.UList("attributes"),
	}

	settings.colors.name = ymlConfig.UString("colors.name", "blue")
	settings.colors.value = ymlConfig.UString("colors.value", "white")

	settings.SetDocumentationPath("pokemon")

	return &settings
}
