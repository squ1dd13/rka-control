package rka

import (
	"bytes"
	"encoding/hex"
	"log"
	"math"
)

type LED struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

func HexToLED(s string) LED {
	if len(s) != 6 && len(s) != 8 {
		log.Fatalf("Invalid colour: %s (expected 'RRGGBB' or 'RRGGBBAA')\n", s)
	}

	// Assume max alpha if nothing specified.
	if len(s) == 6 {
		s += "FF"
	}

	components, err := hex.DecodeString(s)

	if err != nil {
		log.Fatal(err)
	}

	return LED{
		Red:   components[0],
		Green: components[1],
		Blue:  components[2],
		Alpha: components[3],
	}
}

func NewLED(r uint8, g uint8, b uint8) LED {
	return LED{
		r,
		g,
		b,
		0xFF,
	}
}

func RGB(rgb uint32) LED {
	return NewLED(
		uint8((rgb>>16)&255),
		uint8((rgb>>8)&255),
		uint8(rgb&255),
	)
}

func RGBA(r uint8, g uint8, b uint8, a uint8) LED {
	return LED{
		r,
		g,
		b,
		a,
	}
}

func lerp(a uint8, b uint8, t float32) uint8 {
	// Fairly ugly, so this is in a separate function.
	return uint8(math.Round(float64((1-t)*float32(a) + t*float32(b))))
}

func (led *LED) lerp(other *LED, t float32) LED {
	return LED{
		Red:   lerp(led.Red, other.Red, t),
		Green: lerp(led.Green, other.Green, t),
		Blue:  lerp(led.Blue, other.Blue, t),
		Alpha: lerp(led.Alpha, other.Alpha, t),
	}
}

func (led *LED) Write(buf *bytes.Buffer) {
	buf.Write([]byte{
		led.Alpha,
		led.Red,
		led.Green,
		led.Blue,
		0, /* Padding */
	})
}
