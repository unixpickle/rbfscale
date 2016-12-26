package rbfscale

import (
	"math"
	"math/rand"
	"testing"

	"github.com/unixpickle/num-analysis/linalg"
)

func TestRBFMatrixApply(t *testing.T) {
	const width = 25
	const height = 10
	const pixelCount = width * height
	const variance = 0.7
	distMat := linalg.NewMatrix(pixelCount, pixelCount)
	for x1 := 0; x1 < width; x1++ {
		for x2 := 0; x2 < width; x2++ {
			for y1 := 0; y1 < height; y1++ {
				for y2 := 0; y2 < height; y2++ {
					dist := float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))
					rbf := math.Exp(-dist / (2 * variance))
					distMat.Set(x1+y1*width, x2+y2*width, rbf)
				}
			}
		}
	}
	inVec := make(linalg.Vector, pixelCount)
	for i := range inVec {
		inVec[i] = rand.NormFloat64()
	}
	expected := distMat.Mul(linalg.NewMatrixColumn(inVec)).Data
	testMat := &rbfMatrix{
		Variance: variance,
		Width:    width,
		Height:   height,
	}
	actual := testMat.Apply(inVec)
	if actual.Copy().Scale(-1).Add(expected).MaxAbs() > 1e-5 {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func BenchmarkRBFMatrixApply(b *testing.B) {
	mat := &rbfMatrix{
		Variance: 3,
		Width:    30,
		Height:   30,
	}
	inVec := make(linalg.Vector, mat.Width*mat.Height)
	for i := range inVec {
		inVec[i] = rand.NormFloat64()
	}
	// Fill whatever caches we might use.
	mat.Apply(inVec)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mat.Apply(inVec)
	}
}
