package main

import (
	"reflect"
	"strings"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// wanna do themes now, can fix all the overlay stuff later.
// just some basic io/reading

// not sure what i should use, i guess lua for themes aswell would be nice
// than you can automate a lot of stuff

type Style struct {
	FG           tcell.Color
	BG           tcell.Color
	Underline    bool
	StrikeTrough bool
	Bold         bool
	Italic       bool
	Reverse      bool
}

type ThemeCursor struct {
	Main     Style
	Multi    Style
	Disabled Style
}

type Theme struct {
	Cursor    ThemeCursor
	StatusBar Style
	Text      Style
}

func (t *Theme) SetStyle(key, fg, bg, attr string) {
	style := Style{
		FG: tcell.GetColor(fg),
		BG: tcell.GetColor(bg),
	}
	switch attr {
	case "bold":
		style.Bold = true
	case "italic":
		style.Italic = true
	case "underline":
		style.Underline = true
	case "reverse":
		style.Reverse = true
	case "":
		// no-op
	default:
		ThrowErrF("unknown style attribute: %s", attr)
	}

	parts := strings.Split(key, ".")
	val := reflect.ValueOf(t).Elem()

	for i, part := range parts {
		part = cases.Title(language.English).String(strings.ReplaceAll(part, "-", ""))
		field := val.FieldByNameFunc(func(n string) bool {
			return strings.EqualFold(n, part)
		})

		if !field.IsValid() {
			ThrowErrF("invalid theme style key: %s (at part %q)", key, part)
		}

		if i == len(parts)-1 {
			// Final part — must be of type Style
			if field.Type() != reflect.TypeOf(Style{}) {
				ThrowErrF("field %s is not a Style", part)
			}
			field.Set(reflect.ValueOf(style))
		} else {
			// Must be a nested struct
			if field.Kind() == reflect.Struct {
				val = field
			} else {
				ThrowErrF("cannot descend into non-struct field: %s", part)
			}
		}
	}
}

func (s *Style) Darken(factor float32) Style {
	if s == nil {
		return Style{}
	}
	if factor == 0 {
		return *s
	}
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}

	darkened := *s

	if s.FG.TrueColor().Hex() != -1 {
		rf, gf, bf := s.FG.RGB()
		darkened.FG = tcell.NewRGBColor(
			int32(float32(rf)*(1-factor)),
			int32(float32(gf)*(1-factor)),
			int32(float32(bf)*(1-factor)),
		)
	}

	if s.BG.TrueColor().Hex() != -1 {
		rb, gb, bb := s.BG.RGB()
		darkened.BG = tcell.NewRGBColor(
			int32(float32(rb)*(1-factor)),
			int32(float32(gb)*(1-factor)),
			int32(float32(bb)*(1-factor)),
		)
	}

	return darkened
}

func (s *Style) TCellStyle() tcell.Style {
	var style tcell.Style = tcell.StyleDefault

	style = style.Foreground(s.FG)
	style = style.Background(s.BG)
	style = style.Bold(s.Bold)
	style = style.Italic(s.Italic)
	style = style.Underline(s.Underline)
	style = style.StrikeThrough(s.StrikeTrough)
	style = style.Reverse(s.Reverse)

	return style
}

// func ColorOpacity(color tcell.Color, alpha float32) tcell.Color {
// 	hex := color.TrueColor().String()
// 	if hex == "" {
// 		return color // No alpha support for non-true colors
// 	}
// 	hex = strings.TrimPrefix(hex, "#")
// 	if len(hex) != 6 && len(hex) != 8 {
// 		ThrowErrF("invalid color hex: %s", hex)
// 	}
// 	opacity := int(alpha * 255)
// 	fmt.Println(opacity)
// 	if len(hex) == 6 {
// 		hex += "FF" // Add full opacity if missing
// 	} else {
// 		hex = hex[:6] + "FF" // Ensure we have 8 characters
// 	}

// 	return tcell.Color100
// 	// Convert to RGBA
// }

// func BlendColor(src, dst tcell.Color, alpha float64) tcell.Color {
// 	// Clamp alpha to 0–1
// 	if alpha < 0 {
// 		alpha = 0
// 	} else if alpha > 1 {
// 		alpha = 1
// 	}

// 	sr, sg, sb := src.RGB()
// 	dr, dg, db := dst.RGB()

// 	// tcell gives 0–255 values, we blend them
// 	r := uint8(float64(sr)*(1-alpha) + float64(dr)*alpha)
// 	g := uint8(float64(sg)*(1-alpha) + float64(dg)*alpha)
// 	b := uint8(float64(sb)*(1-alpha) + float64(db)*alpha)

// 	return tcell.NewRGBColor(int32(r), int32(g), int32(b))
// }
