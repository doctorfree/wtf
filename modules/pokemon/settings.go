package pokemon

import (
	"time"

	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "Pokémon"
)

type colors struct {
	name         string
	random_name  string
	value        string
	random_value string
}

type Settings struct {
	*cfg.Common

	colors
	pokemon_en     string
	pokemon_id     int
	pokemon_name   string
	random         bool
	language       string
	staticInterval time.Duration
	randomInterval time.Duration
	requestTimeout int
	attributes     []interface{} `help:"Defines what data to display and the order." values:"'size', 'experience', 'genus', and/or 'text'"`
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		pokemon_en:     ymlConfig.UString("pokemon_name", ""),
		pokemon_id:     ymlConfig.UInt("pokemon_id", 0),
		pokemon_name:   ymlConfig.UString("pokemon_name", ""),
		random:         ymlConfig.UBool("random", true),
		language:       ymlConfig.UString("language", "en"),
		staticInterval: cfg.ParseTimeString(ymlConfig, "staticInterval", "300s"),
		randomInterval: cfg.ParseTimeString(ymlConfig, "randomInterval", "30s"),
		requestTimeout: ymlConfig.UInt("timeout", 30),
		attributes:     ymlConfig.UList("attributes"),
	}

	settings.colors.name = ymlConfig.UString("colors.name", "blue")
	settings.colors.random_name = ymlConfig.UString("colors.random_name", "lightblue")
	settings.colors.value = ymlConfig.UString("colors.value", "yellow")
	settings.colors.random_value = ymlConfig.UString("colors.random_value", "cyan")

	if settings.random {
		settings.RefreshInterval = settings.randomInterval
	} else {
		settings.RefreshInterval = settings.staticInterval
	}

	settings.SetDocumentationPath("pokemon")

	return &settings
}
