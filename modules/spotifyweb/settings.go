package spotifyweb

import (
	"os"

	"github.com/olebedev/config"
	"github.com/doctorfree/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "Spotify Web"
)

type Settings struct {
	*cfg.Common

	callbackPort string
	clientID     string
	secretKey    string
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {

	settings := Settings{
		Common: cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),

		callbackPort: ymlConfig.UString("callbackPort", "8080"),
		clientID:     ymlConfig.UString("clientID", os.Getenv("SPOTIFY_ID")),
		secretKey:    ymlConfig.UString("secretKey", os.Getenv("SPOTIFY_SECRET")),
	}

	cfg.ModuleSecret(name, globalConfig, &settings.secretKey).Load()

	return &settings
}
