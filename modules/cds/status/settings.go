package cdsstatus

import (
	"os"

	"github.com/olebedev/config"
	"github.com/doctorfree/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "CDS Status"
)

// Settings defines the configuration properties for this module
type Settings struct {
	*cfg.Common

	token  string `help:"Your CDS API token."`
	apiURL string `help:"Your CDS API URL."`
	uiURL  string
}

// NewSettingsFromYAML creates a new settings instance from a YAML config block
func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {
	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		token:  ymlConfig.UString("token", ymlConfig.UString("token", os.Getenv("CDS_TOKEN"))),
		apiURL: ymlConfig.UString("apiURL", os.Getenv("CDS_API_URL")),
	}

	settings.SetDocumentationPath("cds/status")

	return &settings
}
