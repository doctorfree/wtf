package pokemon

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/rivo/tview"
	"github.com/doctorfree/wtf/view"
	"github.com/doctorfree/wtf/wtf"
)

type Widget struct {
	view.TextWidget

	result   string
	settings *Settings
}

func NewWidget(tviewApp *tview.Application, redrawChan chan bool, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, nil, settings.Common),

		settings: settings,
	}

	return &widget
}

func (widget *Widget) Refresh() {
	widget.pokemon()

	widget.Redraw(func() (string, string, bool) { return widget.CommonSettings().Title, widget.result, false })
}

// this method reads the config and calls pokeapi.co for lunar phase
func (widget *Widget) pokemon() {
	client := &http.Client{}

	id := widget.settings.id
	if widget.settings.random {
		rand.Seed(time.Now().UnixNano())
		id = rand.Intn(905) + 1
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", id), http.NoBody)
	if err != nil {
		widget.result = err.Error()
		return
	}

	req.Header.Set("User-Agent", "curl")
	response, err := client.Do(req)
	if err != nil {
		widget.result = err.Error()
		return

	}
	defer func() { _ = response.Body.Close() }()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		widget.result = err.Error()
		return
	}

	widget.result = strings.TrimSpace(wtf.ASCIItoTviewColors(string(contents)))
}
