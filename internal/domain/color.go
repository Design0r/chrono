package domain

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type color struct{}

var Color = color{}

func (color) HSL(seed int) string {
	hue := (seed * 12345) % 360
	saturation := 50
	lightness := 40
	return fmt.Sprintf("hsl(%d, %d%%, %d%%)", hue, saturation, lightness)
}

func (color) HSLDark(seed int) string {
	hue := (seed * 12345) % 360
	saturation := 30
	lightness := 30
	return fmt.Sprintf("hsl(%d, %d%%, %d%%)", hue, saturation, lightness)
}

func (color) HSLFloat(seed int) (float64, float64, float64) {
	hue := float64((seed * 12345) % 360)
	saturation := 0.5
	lightness := 0.5
	return hue, saturation, lightness
}

func (color) RandomHSL() (float64, float64, float64) {
	hue := float64((time.Now().UnixNano() * 12345) % 360)
	saturation := 0.5
	lightness := 0.5

	return hue, saturation, lightness
}

func (color) HSLDarkFromHex(hex string) string {
	h, _, _ := Color.HexToHSL(hex)
	s := 0.3
	l := 0.3
	return Color.HSLToString(h, s, l)
}

func (color) RandomHexColor() string {
	return Color.HSLToHex(Color.RandomHSL())
}

func (color) HSLToHex(h float64, s float64, l float64) string {
	// 1. Normalize the hue to [0..1]
	//    H (in degrees) / 360 => range [0..1]
	h = math.Mod(h, 360) / 360

	// 2. If there is no saturation, the color is a shade of gray
	//    so R=G=B=L
	var r, g, b float64
	if s == 0 {
		r, g, b = l, l, l
	} else {
		// For HSL-to-RGB, we define some helpers:
		var hueToRGB = func(p, q, t float64) float64 {
			// Wrap t around if it goes out of bounds
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			switch {
			case t < 1.0/6.0:
				return p + (q-p)*6.0*t
			case t < 1.0/2.0:
				return q
			case t < 2.0/3.0:
				return p + (q-p)*(2.0/3.0-t)*6.0
			default:
				return p
			}
		}

		// q and p are temporary variables used to compute the intermediate values
		// based on whether Lightness < 0.5 or not.
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		// Convert h, s, l to r, g, b in [0..1]
		r = hueToRGB(p, q, h+1.0/3.0)
		g = hueToRGB(p, q, h)
		b = hueToRGB(p, q, h-1.0/3.0)
	}

	// 3. Scale up to [0..255] and format as hex
	R := int(math.Round(r * 255.0))
	G := int(math.Round(g * 255.0))
	B := int(math.Round(b * 255.0))

	return fmt.Sprintf("#%02X%02X%02X", R, G, B)
}

func (color) HexToHSL(hex string) (float64, float64, float64) {
	// Remove leading "#" if present.
	hex = strings.TrimPrefix(hex, "#")

	// Check for valid length (must be 6 characters: RRGGBB).
	if len(hex) != 6 {
		return 0, 0, 0
	}

	// Parse the hex string as a 24-bit integer.
	rgbValue, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, 0, 0
	}

	// Extract the red, green, and blue components (0–255).
	r := float64((rgbValue>>16)&0xFF) / 255.0
	g := float64((rgbValue>>8)&0xFF) / 255.0
	b := float64(rgbValue&0xFF) / 255.0

	// Find min and max values of R, G, B.
	maxVal := math.Max(r, math.Max(g, b))
	minVal := math.Min(r, math.Min(g, b))
	l := (maxVal + minVal) / 2.0

	var h, s float64

	if maxVal == minVal {
		// Achromatic (grey); no saturation.
		h = 0
		s = 0
	} else {
		// Chromatic colors
		d := maxVal - minVal

		// Saturation
		if l > 0.5 {
			s = d / (2.0 - maxVal - minVal)
		} else {
			s = d / (maxVal + minVal)
		}

		// Hue
		switch maxVal {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h /= 6.0
	}

	// Convert hue to degrees [0..360].
	hDegrees := h * 360.0

	return hDegrees, s, l
}

func (color) HSLToString(h float64, s float64, l float64) string {
	return fmt.Sprintf("hsl(%v, %v%%, %v%%)", h, s*100, l*100)
}
