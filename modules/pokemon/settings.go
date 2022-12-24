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
	name     string
	random   bool
	language string
	attributes []interface{} `help:"Defines what data to display and the order." values:"'height', 'weight', 'name', 'genus', and/or 'id'"`
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		id:         ymlConfig.UInt("id", 0),
		name:       ymlConfig.UString("name", ""),
		random:     ymlConfig.UBool("random", true),
		language:   ymlConfig.UString("language", "en"),
		attributes: ymlConfig.UList("attributes"),
	}

	settings.colors.name = ymlConfig.UString("colors.name", "cyan")
	settings.colors.value = ymlConfig.UString("colors.value", "white")

	settings.SetDocumentationPath("pokemon")

	return &settings
}
