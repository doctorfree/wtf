package security

import (
	"github.com/olebedev/config"
	"github.com/doctorfree/wtf/cfg"
)

const (
	defaultFocusable = false
	defaultTitle     = "Security"
)

type Settings struct {
	*cfg.Common
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),
	}

	return &settings
}
