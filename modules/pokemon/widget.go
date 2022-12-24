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

var attLookup = map[string]string {
	"pokemon_id":    "Species ID",
	"pokemon_name":  "Species Name",
	"genus":         "Species Genus",
	"height":        "Height",
	"weight":        "Weight",
	"text":          "",
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

func (widget *Widget) pokemon() {
	client := &http.Client{}

	id_config := widget.settings.pokemon_id
	name_config := widget.settings.pokemon_name

	if widget.settings.random {
		name_config = ""
		rand.Seed(time.Now().UnixNano())
		id_config = rand.Intn(905) + 1
	}
	idstr := strconv.Itoa(id_config)

	qstr := ""
	if name_config == "" {
		qstr = idstr
	} else {
		qstr = name_config
		qstr = strings.ToLower(name_config)
	}
	req, err := http.NewRequest("GET", "https://pokeapi.co/api/v2/pokemon/"+qstr, http.NoBody)
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

	attrs := utils.ToStrs(widget.settings.attributes)

	if len(attrs) == 0 {
		attrs = []string{"id", "name", "height", "weight", "genus", "text"}
	}

	format := ""

	for _, att := range attrs {
		if val, ok := attLookup[strings.ToLower(att)]; ok {
			format = format + formatableText(val, strings.ToLower(att))
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

	langconfig = widget.settings.language
	pokemon_text := "Unknown"
	for i := range spec.FlavorTextEntries {
        if spec.FlavorTextEntries[i].Language.Name == langconfig {
			pokemon_text = spec.FlavorTextEntries[i].FlavorText
        }
    }
	if pokemon_text == "Unknown" {
		langconfig = "en"
	    for i := range spec.FlavorTextEntries {
            if spec.FlavorTextEntries[i].Language.Name == langconfig {
			    pokemon_text = spec.FlavorTextEntries[i].FlavorText
            }
        }
	}

	err := resultTemplate.Execute(resultBuffer, map[string]string{
		"nameColor":     widget.settings.colors.name,
		"valueColor":    widget.settings.colors.value,
		"pokemon_id":    strconv.Itoa(spec.ID),
		"pokemon_name":  pokemon_name,
		"genus":         pokemon_genus,
		"height":        fmt.Sprintf("%6.2f (m)", float64(poke.Height) / 10.0),
		"weight":        fmt.Sprintf("%6.2f (kg)", float64(poke.Weight) / 10.0),
		"text":          pokemon_text,
	})

	if err != nil {
		widget.result = err.Error()
	}

	widget.result = resultBuffer.String()
}

func formatableText(key, value string) string {
	return fmt.Sprintf(" [{{.nameColor}}]%s: [{{.valueColor}}]{{.%s}}\n", key, value)
}
