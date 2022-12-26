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
	view.ScrollableWidget

	result   string
	settings *Settings
}

var attLookup = map[string]string {
	"pokemon_name":  "Species Name",
	"genus":         "Species Genus",
	"pokemon_id":    "Species ID",
	"pokemon_types": "Pokémon Types",
	"experience":    "Base Experience",
	"size":          "Height (m) / Weight (kg)",
	"text":          "Description",
}

func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := &Widget{
		ScrollableWidget: view.NewScrollableWidget(tviewApp, redrawChan, pages, settings.Common),
		settings:         settings,
	}

	widget.SetRenderFunction(widget.Refresh)
	widget.initializeKeyboardControls()

	return widget
}

func (widget *Widget) Refresh() {
	widget.pokemon()

	if !widget.settings.Enabled {
		widget.View.Clear()
		return
	}

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
		widget.settings.RefreshInterval = widget.settings.randomInterval
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
		attrs = []string{"pokemon_id", "pokemon_name", "size", "genus", "pokemon_types", "experience", "text"}
	}

	format := ""
	attlower := ""
	for _, att := range attrs {
		attlower = strings.ToLower(att)
		if val, ok := attLookup[attlower]; ok {
			format = format + formatableText(val, attlower)
		}
	}

	resultTemplate, _ := template.New("pokemon_result").Parse(format)

	langconfig := widget.settings.language
	resultBuffer := new(bytes.Buffer)
	pokemon_name := "Unknown"
	en_pokemon_name := "Unknown"
	for i := range spec.Names {
        if spec.Names[i].Language.Name == langconfig {
			pokemon_name = spec.Names[i].Name
        }
        if spec.Names[i].Language.Name == "en" {
			en_pokemon_name = spec.Names[i].Name
        }
    }
	if pokemon_name == "Unknown" {
		langconfig = "en"
	    for i := range spec.Names {
            if spec.Names[i].Language.Name == langconfig {
			    pokemon_name = spec.Names[i].Name
				en_pokemon_name = pokemon_name
            }
        }
	}
	widget.settings.pokemon_en = en_pokemon_name
	widget.settings.pokemon_id = spec.ID

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
	pokemon_text := "\nUnknown"
	for i := range spec.FlavorTextEntries {
        if spec.FlavorTextEntries[i].Language.Name == langconfig {
			pokemon_text = "\n" + spec.FlavorTextEntries[i].FlavorText
        }
    }
	if pokemon_text == "\nUnknown" {
		langconfig = "en"
	    for i := range spec.FlavorTextEntries {
            if spec.FlavorTextEntries[i].Language.Name == langconfig {
			    pokemon_text = "\n" + spec.FlavorTextEntries[i].FlavorText
            }
        }
	}

//	colorReset := "\033[0m"
//  colorRed := "\033[31m"
//  colorGreen := "\033[32m"
//  colorYellow := "\033[33m"
//  colorBlue := "\033[34m"
//  colorPurple := "\033[35m"
//  colorCyan := "\033[36m"
//  colorWhite := "\033[37m"

	pokemon_types := ""
	poketype := ""
//	color := 7
	first := false
	for i := range poke.Types {
		poketype = strings.ToUpper(poke.Types[i].Type.Name)
//		switch poketype {
//		case "NORMAL":
//			color=7
//		case "FIRE":
//			color=9
//		case "WATER":
//			color=12
//		case "ELECTRIC":
//			color=11
//		case "GRASS":
//			color=10
//		case "ICE":
//			color=14
//		case "FIGHTING":
//			color=1
//		case "POISON":
//			color=5
//		case "GROUND":
//			color=11
//		case "FLYING":
//			color=6
//		case "PSYCHIC":
//			color=13
//		case "BUG":
//			color=2
//		case "ROCK":
//			color=3
//		case "GHOST":
//			color=4
//		case "DRAGON":
//			color=4
//		case "DARK":
//			color=3
//		case "STEEL":
//			color=8
//		case "FAIRY":
//			color=13
//		default:
//			color=7
//		}

		if first {
			pokemon_types += " "
		} else {
			first = true
		}
//		pokemon_types += "[7;38;5;" + strconv.Itoa(color) + "m" + poketype + "[0m"
		pokemon_types += poketype
    }

	err := resultTemplate.Execute(resultBuffer, map[string]string{
		"nameColor":      widget.settings.colors.name,
		"valueColor":     widget.settings.colors.value,
		"pokemon_name":   pokemon_name,
		"genus":          pokemon_genus,
		"pokemon_id":     strconv.Itoa(spec.ID),
		"pokemon_types":  pokemon_types,
		"experience":     strconv.Itoa(poke.BaseExperience),
		"size":           fmt.Sprintf("%.2f / %.2f", float64(poke.Height) / 10.0, float64(poke.Weight) / 10.0),
		"text":           pokemon_text,
	})

	if err != nil {
		widget.result = err.Error()
	}

	widget.result = resultBuffer.String()
}

func formatableText(key, value string) string {
	return fmt.Sprintf(" [{{.nameColor}}]%s: [{{.valueColor}}]{{.%s}}\n", key, value)
}

// NextPokemon shows the next Pokemon ID or wraps to the ID 1
func (widget *Widget) NextPokemon() {

	if widget.settings.random {
		return
	}

	curr_id := widget.settings.pokemon_id

	if curr_id >= 905 {
		widget.settings.pokemon_id = 1
	} else {
		widget.settings.pokemon_id = curr_id + 1
	}
	widget.settings.pokemon_name = ""
	widget.Refresh()
}

// PrevPokemon shows the previous Pokemon ID or wraps to ID 905
func (widget *Widget) PrevPokemon() {

	if widget.settings.random {
		return
	}

	curr_id := widget.settings.pokemon_id

	if curr_id == 1 {
		widget.settings.pokemon_id = 905
	} else {
		widget.settings.pokemon_id = curr_id - 1
	}
	widget.settings.pokemon_name = ""
	widget.Refresh()
}

// https://bulbapedia.bulbagarden.net/wiki/Bulbasaur_(Pok%C3%A9mon)
func (widget *Widget) OpenPokemon() {
	poke_name := widget.settings.pokemon_en
	if poke_name == "Unknown" {
		return
	}
	if poke_name == "" {
		return
	}
	poke_url := "https://bulbapedia.bulbagarden.net/wiki/" + poke_name + "_(Pok%C3%A9mon)"
	utils.OpenFile(poke_url)
}

// ToggleRandom toggles the random display of Pokemon
func (widget *Widget) ToggleRandom() {

	if widget.settings.random {
		widget.settings.random = false
		// Restore refreshInterval config when in static mode
		widget.settings.RefreshInterval = widget.settings.interval
	} else {
		widget.settings.random = true
		// Ignore refreshInterval config when in random mode
		// TODO: test if saved 'interval' is less than 60s
		widget.settings.RefreshInterval = widget.settings.randomInterval
	}

	widget.settings.pokemon_name = ""
	widget.Refresh()
}

// Disable/Enable the widget
func (widget *Widget) DisableWidget() {

	if widget.settings.Enabled {
		widget.settings.Enabled = false
	} else {
		widget.settings.Enabled = true
	}

	widget.Refresh()
}
