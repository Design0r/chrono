package domain_test

import (
	"math"
	"strings"
	"testing"

	"chrono/internal/domain"
)

func TestRandomHexColor(t *testing.T) {
	color := domain.RandomHexColor()
	if !strings.HasPrefix(color, "#") {
		t.Errorf("Expected color to start with '#', got %s", color)
	}
	if len(color) != 7 {
		t.Errorf("Expected hex color length of 7 (like #FFFFFF), got %d", len(color))
	}
}

func TestHSLToHex(t *testing.T) {
	tests := []struct {
		h, s, l float64
	}{
		{0, 1, 0.5},     // pure red
		{120, 1, 0.5},   // pure green
		{240, 1, 0.5},   // pure blue
		{360, 1, 0.5},   // same as 0
		{180, 0.5, 0.5}, // some teal color
		{359, 0.0, 0.5}, // grey, since saturation=0
	}
	for _, tc := range tests {
		got := domain.HSLToHex(tc.h, tc.s, tc.l)
		if !strings.HasPrefix(got, "#") || len(got) != 7 {
			t.Errorf("HSLToHex(%v, %v, %v) = %s, invalid format",
				tc.h, tc.s, tc.l, got)
		}
	}
}

func TestHexToHSL(t *testing.T) {
	// #FF0000 => h=0, s=1, l=0.5
	h, s, l := domain.HexToHSL("#FF0000")
	if math.Abs(h-0) > 0.1 || math.Abs(s-1) > 0.01 || math.Abs(l-0.5) > 0.01 {
		t.Errorf("Expected #FF0000 => h=0, s=1, l=0.5, got h=%v, s=%v, l=%v", h, s, l)
	}
}
