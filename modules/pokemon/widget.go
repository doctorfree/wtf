package pokemon

import (
	"bytes"
    "encoding/json"
    "fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/rivo/tview"
	"github.com/doctorfree/wtf/utils"
	"github.com/doctorfree/wtf/view"
)

type Widget struct {
	view.TextWidget

	result   string
	settings *Settings
}

var argLookup = map[string]string {
	"id":            "Species ID",
	"name":          "Species Name",
	"genus":         "Species Genus",
	"height":        "Height (m)",
	"weight":        "Weight (kg)",
}

func NewWidget(tviewApp *tview.Application, redrawChan chan bool, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, nil, settings.Common),

		settings: settings,
	}

	widget.View.SetWrap(false)

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

	if err != nil {
		widget.result = err.Error()
		return
	}
	defer func() { _ = response.Body.Close() }()

    var pokemonObject Pokemon
	err = json.NewDecoder(response.Body).Decode(&pokemonObject)
	if err != nil {
		widget.result = err.Error()
		return
	}

	spreq, err := http.NewRequest("GET", pokemonObject.Species.URL, http.NoBody)
	if err != nil {
		widget.result = err.Error()
		return
	}

	spreq.Header.Set("User-Agent", "curl")
	spresponse, err := client.Do(spreq)

	if err != nil {
		widget.result = err.Error()
		return
	}
	defer func() { _ = spresponse.Body.Close() }()

    var speciesObject PokemonSpecies
	err = json.NewDecoder(spresponse.Body).Decode(&speciesObject)
	if err != nil {
		widget.result = err.Error()
		return
	}

	widget.setResult(&pokemonObject, &speciesObject)
}

func (widget *Widget) setResult(poke *Pokemon, spec *PokemonSpecies) {

	args := utils.ToStrs(widget.settings.args)

	// if no arguments are defined set default
	if len(args) == 0 {
		args = []string{"id", "name", "height", "weight", "genus"}
	}

	format := ""

	for _, arg := range args {
		if val, ok := argLookup[strings.ToLower(arg)]; ok {
			format = format + formatableText(val, strings.ToLower(arg))
		}
	}

	resultTemplate, _ := template.New("pokemon_result").Parse(format)

	langconfig := widget.settings.language
	resultBuffer := new(bytes.Buffer)
	pokemon_name := "Unknown"
	for i := range spec.Names {
        if spec.Names[i].Language.Name == langconfig {
			pokemon_name = spec.Names[i].Name
        }
    }
	if pokemon_name == "Unknown" {
		langconfig = "en"
	    for i := range spec.Names {
            if spec.Names[i].Language.Name == langconfig {
			    pokemon_name = spec.Names[i].Name
            }
        }
	}

	langconfig = widget.settings.language
	pokemon_genus := "Unknown"
	for i := range spec.Genera {
        if spec.Genera[i].Language.Name == langconfig {
			pokemon_genus = spec.Genera[i].Genus
        }
    }
	if pokemon_genus == "Unknown" {
		langconfig = "en"
	    for i := range spec.Genera {
            if spec.Genera[i].Language.Name == langconfig {
			    pokemon_genus = spec.Genera[i].Genus
            }
        }
	}

	err := resultTemplate.Execute(resultBuffer, map[string]string{
		"nameColor":     widget.settings.colors.name,
		"valueColor":    widget.settings.colors.value,
		"id":            strconv.Itoa(spec.ID),
		"name":          pokemon_name,
		"genus":         pokemon_genus,
		"height":        strconv.Itoa(poke.Height),
		"weight":        strconv.Itoa(poke.Weight),
	})

	if err != nil {
		widget.result = err.Error()
	}

//	idstr := strconv.Itoa(pokemon_id)
//  fmt.Println("▐[1;7m No. " + idstr + "[0m▌ [1m" + pokemon_name + " - " + pokemon_genus + "[0m")

	widget.result = resultBuffer.String()
}

func formatableText(key, value string) string {
	return fmt.Sprintf(" [{{.nameColor}}]%s: [{{.valueColor}}]{{.%s}}\n", key, value)
}
