package pokemon

import "github.com/gdamore/tcell/v2"

func (widget *Widget) initializeKeyboardControls() {
	widget.InitializeHelpTextKeyboardControl(widget.ShowHelp)
	widget.InitializeRefreshKeyboardControl(widget.Refresh)

	widget.SetKeyboardChar("n", widget.NextPokemon, "Select next Pokémon")
	widget.SetKeyboardChar("p", widget.PrevPokemon, "Select previous Pokémon")
	widget.SetKeyboardChar("o", widget.OpenPokemon, "Open Pokémon at Bulbapedia in browser")
	widget.SetKeyboardChar("R", widget.ToggleRandom, "Toggle random Pokémon display")

	widget.SetKeyboardKey(tcell.KeyLeft, widget.PrevPokemon, "Select previous Pokémon")
	widget.SetKeyboardKey(tcell.KeyRight, widget.NextPokemon, "Select next Pokémon")
	widget.SetKeyboardKey(tcell.KeyEnter, widget.OpenPokemon, "Open Pokémon at Bulbapedia in browser")
	widget.SetKeyboardKey(tcell.KeyUp, widget.ToggleRandom, "Toggle random Pokémon display")
	widget.SetKeyboardKey(tcell.KeyCtrlD, widget.DisableWidget, "Disable/Enable this Pokémon widget instance")
}
