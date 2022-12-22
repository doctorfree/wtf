package nbascore

import (
	"github.com/olebedev/config"
	"github.com/doctorfree/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "NBA Score"
)

type Settings struct {
	*cfg.Common
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),
	}

	settings.SetDocumentationPath("sports/nbascore")

	return &settings
}
