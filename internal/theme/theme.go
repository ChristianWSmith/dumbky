package theme

import (
	"dumbky/internal/log"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type DumbkyTheme struct{}

func (DumbkyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	// Colors
	case theme.ColorNameHyperlink:
		return color.NRGBA{R: 246, G: 188, B: 40, A: 255} // #f6bc28
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 246, G: 188, B: 40, A: 255} // #f6bc28
	case theme.ColorNameSelection:
		return color.NRGBA{R: 255, G: 183, B: 0, A: 64} // #ffb700
	case theme.ColorNameFocus:
		return color.NRGBA{R: 255, G: 183, B: 0, A: 64} // #ffb700
	case theme.ColorNameError:
		return color.NRGBA{R: 224, G: 67, B: 54, A: 255} // #e04336
	case theme.ColorNameWarning:
		return color.NRGBA{R: 224, G: 67, B: 54, A: 255} // #e04336
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 67, G: 244, B: 54, A: 255} // #43f436

	// Off-Grays
	case theme.ColorNameBackground:
		return color.NRGBA{R: 24, G: 23, B: 23, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 46, G: 41, B: 40, A: 255}
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 46, G: 41, B: 40, A: 255}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 58, G: 57, B: 57, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 35, G: 32, B: 32, A: 255}
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 58, G: 57, B: 57, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 37, G: 29, B: 24, A: 255}
	case theme.ColorNameScrollBarBackground:
		return color.NRGBA{R: 35, G: 32, B: 32, A: 255}
	case theme.ColorNameForegroundOnSuccess:
		return color.NRGBA{R: 24, G: 23, B: 23, A: 255}
	case theme.ColorNameForegroundOnWarning:
		return color.NRGBA{R: 24, G: 23, B: 23, A: 255}
	case theme.ColorNameForegroundOnError:
		return color.NRGBA{R: 24, G: 23, B: 23, A: 255}

	// Grays
	case theme.ColorNameForegroundOnPrimary:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 243, G: 243, B: 243, A: 255}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 27, G: 27, B: 27, A: 255}
	case theme.ColorNameHover:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 15}
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 40, G: 0x29, B: 0x2e, A: 255}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 178, G: 178, B: 178, A: 255}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 102}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 153}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0, G: 0, B: 0, A: 102}

	// Default
	default:
		log.Warn(fmt.Errorf("unexpected color name %s", name))
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
