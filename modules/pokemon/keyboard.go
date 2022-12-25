package pokemon

import "github.com/gdamore/tcell/v2"

func (widget *Widget) initializeKeyboardControls() {
	widget.InitializeHelpTextKeyboardControl(widget.ShowHelp)
	widget.InitializeRefreshKeyboardControl(widget.Refresh)

	widget.SetKeyboardChar("n", widget.NextPokemon, "Select next Pokémon")
	widget.SetKeyboardChar("p", widget.PrevPokemon, "Select previous Pokémon")

	widget.SetKeyboardKey(tcell.KeyLeft, widget.PrevPokemon, "Select previous Pokémon")
	widget.SetKeyboardKey(tcell.KeyRight, widget.NextPokemon, "Select next Pokémon")
}
