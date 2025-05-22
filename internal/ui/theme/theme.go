package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type DumbkyTheme struct{}

func (DumbkyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xe9, G: 0x83, B: 0x19, A: 0x7f}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0xf3, G: 0xb8, B: 0x27, A: 0xff}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (DumbkyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (DumbkyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (DumbkyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
