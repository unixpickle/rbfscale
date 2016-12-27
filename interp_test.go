package rbfscale

import (
	"image"
	"image/color"
	"math/rand"
	"testing"
)

func BenchmarkNewInterp(b *testing.B) {
	const width = 8
	const height = 10
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			img.Set(j, i, color.RGBA{
				R: uint8(rand.Intn(0x100)),
				G: uint8(rand.Intn(0x100)),
				B: uint8(rand.Intn(0x100)),
				A: 0xff,
			})
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewInterp(img, 2)
	}
}
