package pokemon

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

//	"github.com/mtslzr/pokeapi-go"
//	"github.com/mtslzr/pokeapi-go/structs"
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

// this method reads the config and calls pokeapi.co for the Pokemon ID
func (widget *Widget) pokemon() {
	client := &http.Client{}

	id := widget.settings.id
	if widget.settings.random {
		rand.Seed(time.Now().UnixNano())
		id = rand.Intn(905) + 1
	}

	idstr := strconv.Itoa(id)
	req, err := http.NewRequest("GET", "https://pokeapi.co/api/v2/pokemon/"+idstr, http.NoBody)
	if err != nil {
		widget.result = err.Error()
		return
	}

	req.Header.Set("User-Agent", "curl")
	response, err := client.Do(req)
//	var response structs.Pokemon

//	response, err := pokeapi.Pokemon(strconv.Itoa(id))
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
