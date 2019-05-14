package blurhash

import (
	"fmt"
	"math"
	"testing"
)

func TestSignPow(t *testing.T) {
	var testCases = []struct {
		in  float64
		exp float64
		out float64
	}{
		{-2.0, 4, -16.0},
		{2.0, 4, 16.0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%f", tc.in), func(t *testing.T) {
			sp := signPow(tc.in, tc.exp)
			if sp != tc.out {
				t.Fatalf("got %f, wanted %f", sp, tc.out)
			}
		})
	}
}

func TestLinearTosRGB(t *testing.T) {
	var testCases = []struct {
		in  float64
		out int
	}{
		{0.0, 0},
		{1.0, 255},
		{0.5, 188},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%f", tc.in), func(t *testing.T) {
			sp := linearTosRGB(tc.in)
			if sp != tc.out {
				t.Fatalf("got %d, wanted %d", sp, tc.out)
			}
		})
	}
}

func TestSRGBToLinear(t *testing.T) {
	var testCases = []struct {
		in  int
		out float64
	}{
		{0, 0.0},
		{255, 1.0},
		{188, 0.5},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d", tc.in), func(t *testing.T) {
			sp := sRGBToLinear(tc.in)
			if math.Abs(tc.out-sp) > 0.05 {
				t.Fatalf("got %f, wanted %f", sp, tc.out)
			}
		})
	}
}
