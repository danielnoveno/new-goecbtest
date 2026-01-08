/*
    file:           views/theme/mytheme.go
    description:    Tema UI untuk mytheme
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package theme

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"go-ecb/app/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type Palette struct {
	Name            string
	Description     string
	Background      color.Color
	Foreground      color.Color
	Text            color.Color
	Button          color.Color
	Disabled        color.Color
	Error           color.Color
	Focus           color.Color
	Hover           color.Color
	InputBackground color.Color
	Placeholder     color.Color
	Primary         color.Color
	Scrollbar       color.Color
	Selection       color.Color
	HeaderColor     color.Color
	Navbar          color.Color
	Footer          color.Color
	Accent          color.Color
}

var fallbackPalette = Palette{
	Name:            "Polytron Aurora",
	Description:     "Palet gelap dengan kontras lebih tinggi untuk teks terang di latar putih.",
	Background:      color.RGBA{R: 0x0F, G: 0x17, B: 0x2A, A: 0xFF},
	Foreground:      color.RGBA{R: 0xE2, G: 0xE8, B: 0xF4, A: 0xFF},
	Text:            color.RGBA{R: 0xF8, G: 0xFA, B: 0xFC, A: 0xFF},
	Button:          color.RGBA{R: 0x22, G: 0x8F, B: 0xDD, A: 0xFF},
	Disabled:        color.RGBA{R: 0x4B, G: 0x4F, B: 0x62, A: 0xFF},
	Error:           color.RGBA{R: 0xF8, G: 0x71, B: 0x71, A: 0xFF},
	Focus:           color.RGBA{R: 0x3B, G: 0x82, B: 0xF6, A: 0xFF},
	Hover:           color.RGBA{R: 0x1E, G: 0x24, B: 0x38, A: 0xFF},
	InputBackground: color.RGBA{R: 0x1F, G: 0x29, B: 0x3C, A: 0xFF},
	Placeholder:     color.RGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF},
	Primary:         color.RGBA{R: 0x38, G: 0xB4, B: 0xF5, A: 0xFF},
	Scrollbar:       color.RGBA{R: 0x0F, G: 0x17, B: 0x2A, A: 0xFF},
	Selection:       color.RGBA{R: 0x38, G: 0xB4, B: 0xF5, A: 0xFF},
	HeaderColor:     color.RGBA{R: 0x0C, G: 0x16, B: 0x27, A: 0xFF},
	Navbar:          color.RGBA{R: 0x0C, G: 0x16, B: 0x27, A: 0xFF},
	Footer:          color.RGBA{R: 0x0C, G: 0x16, B: 0x27, A: 0xFF},
	Accent:          color.RGBA{R: 0x38, G: 0xB4, B: 0xF5, A: 0xFF},
}

// DefaultPalette adalah fungsi untuk default palette.
func DefaultPalette() Palette {
	return fallbackPalette
}

// PaletteFromRecord adalah fungsi untuk palette from record.
func PaletteFromRecord(rec types.Theme) Palette {
	return Palette{
		Name:            rec.Nama,
		Description:     rec.Keterangan,
		Background:      safeHexColor(rec.ColorBackground, fallbackPalette.Background),
		Foreground:      safeHexColor(rec.ColorForeground, fallbackPalette.Foreground),
		Text:            safeHexColor(rec.ColorText, fallbackPalette.Text),
		Button:          safeHexColor(rec.ColorButton, fallbackPalette.Button),
		Disabled:        safeHexColor(rec.ColorDisabled, fallbackPalette.Disabled),
		Error:           safeHexColor(rec.ColorError, fallbackPalette.Error),
		Focus:           safeHexColor(rec.ColorFocus, fallbackPalette.Focus),
		Hover:           safeHexColor(rec.ColorHover, fallbackPalette.Hover),
		InputBackground: safeHexColor(rec.ColorInputBackground, fallbackPalette.InputBackground),
		Placeholder:     safeHexColor(rec.ColorPlaceholder, fallbackPalette.Placeholder),
		Primary:         safeHexColor(rec.ColorPrimary, fallbackPalette.Primary),
		Scrollbar:       safeHexColor(rec.ColorScrollbar, fallbackPalette.Scrollbar),
		Selection:       safeHexColor(rec.ColorSelection, fallbackPalette.Selection),
		HeaderColor:     safeHexColor(rec.HeaderStart, fallbackPalette.HeaderColor),
		Navbar:          safeHexColor(rec.ColorNavbar, fallbackPalette.Navbar),
		Footer:          safeHexColor(rec.ColorFooter, fallbackPalette.Footer),
		Accent:          safeHexColor(rec.Accent, fallbackPalette.Accent),
	}
}

// safeHexColor adalah fungsi untuk safe hex color.
func safeHexColor(value string, fallback color.Color) color.Color {
	c, err := parseHexColor(value)
	if err != nil {
		return fallback
	}
	return c
}

// ColorFromHex adalah fungsi untuk color from hex.
func ColorFromHex(value string, fallback color.Color) color.Color {
	return safeHexColor(value, fallback)
}

// parseHexColor adalah fungsi untuk parse hex color.
func parseHexColor(value string) (color.RGBA, error) {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "#")
	if len(value) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color %q", value)
	}

	r, err := strconv.ParseUint(value[0:2], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	g, err := strconv.ParseUint(value[2:4], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}
	b, err := strconv.ParseUint(value[4:6], 16, 8)
	if err != nil {
		return color.RGBA{}, err
	}

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 0xFF,
	}, nil
}

type MyTheme struct {
	palette Palette
}

var _ fyne.Theme = (*MyTheme)(nil)

// New adalah fungsi untuk baru.
func New(p Palette) fyne.Theme {
	return &MyTheme{palette: p}
}

// Color adalah fungsi untuk color.
func (m *MyTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return m.palette.Background
	case theme.ColorNameButton:
		return m.palette.Button
	case theme.ColorNameDisabled, theme.ColorNameDisabledButton:
		return m.palette.Disabled
	case theme.ColorNameError:
		return m.palette.Error
	case theme.ColorNameFocus:
		return m.palette.Focus
	case theme.ColorNameForeground:
		return m.palette.Foreground
	case theme.ColorNameForegroundOnError, theme.ColorNameForegroundOnPrimary,
		theme.ColorNameForegroundOnSuccess, theme.ColorNameForegroundOnWarning:
		return m.palette.Text
	case theme.ColorNameHeaderBackground:
		return m.palette.HeaderColor
	case theme.ColorNameHover:
		return m.palette.Hover
	case theme.ColorNameHyperlink:
		return m.palette.Primary
	case theme.ColorNameInputBackground:
		return m.palette.InputBackground
	case theme.ColorNameInputBorder:
		return m.palette.Hover
	case theme.ColorNameMenuBackground, theme.ColorNameOverlayBackground:
		return m.palette.Background
	case theme.ColorNamePlaceHolder:
		return m.palette.Placeholder
	case theme.ColorNamePressed:
		return m.palette.Hover
	case theme.ColorNamePrimary:
		return m.palette.Primary
	case theme.ColorNameScrollBar:
		return m.palette.Scrollbar
	case theme.ColorNameScrollBarBackground:
		return m.palette.InputBackground
	case theme.ColorNameSelection:
		return m.palette.Selection
	case theme.ColorNameSeparator:
		return m.palette.Hover
	case theme.ColorNameShadow:
		return m.palette.Background
	case theme.ColorNameSuccess:
		return m.palette.Primary
	case theme.ColorNameWarning:
		return m.palette.Error
	default:
		col := theme.DefaultTheme().Color(name, variant)
		if col == nil {
			return m.palette.Text
		}
		return col
	}
}

// Font adalah fungsi untuk font.
func (m *MyTheme) Font(name fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(name)
}

// Icon adalah fungsi untuk icon.
func (m *MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size adalah fungsi untuk size.
func (m *MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
