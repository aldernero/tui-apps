package tui

import "github.com/aldernero/gaul"

var palettes = []gaul.Gradient{
	gaul.NewGradientFromNamed([]string{"cyan", "yellow", "orange", "magenta", "purple"}),
	gaul.NewGradientFromNamed([]string{"cyan", "yellow", "magenta"}),
	gaul.NewGradientFromNamed([]string{"cyan", "magenta"}),
	gaul.NewGradientFromNamed([]string{"blue", "gray", "red"}),
	gaul.NewGradientFromNamed([]string{"orange", "magenta", "purple"}),
	gaul.NewGradientFromNamed([]string{"yellow", "orange", "red"}),
	gaul.NewGradientFromNamed([]string{"cyan", "yellow", "orange", "magenta"}),
}
