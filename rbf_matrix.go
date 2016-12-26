package rbfscale

import (
	"math"

	"github.com/unixpickle/num-analysis/linalg"
)

const epsilon = 1e-8

// A rbfMatrix acts on behalf of an RBF layer with centers
// arranged on the points in an image.
// The Variance is used to divide the argument to exp,
// giving the formula exp(||x-c||/(2*Variance)).
type rbfMatrix struct {
	Variance float64
	Width    int
	Height   int

	rbfCache [][]float64
}

// Dim returns the total number of pixels.
func (n *rbfMatrix) Dim() int {
	return n.Width * n.Height
}

// Apply performs matrix multiplication by the positive
// definite matrix Gij = Rho(Pi-Pj), where Pi is the i-th
// point in the image.
func (n *rbfMatrix) Apply(in linalg.Vector) linalg.Vector {
	result := make(linalg.Vector, 0, n.Dim())
	for y := 0; y < n.Height; y++ {
		for x := 0; x < n.Width; x++ {
			applied := n.applyPoint(x, y, in)
			result = append(result, applied)
		}
	}
	return result
}

func (n *rbfMatrix) applyPoint(x, y int, in linalg.Vector) float64 {
	minX, minY, maxX, maxY := pointBounds(float64(x), float64(y), n.Variance)
	var sum float64
	for i := minY; i < maxY && i < n.Height; i++ {
		for j := minX; j < maxX && j < n.Width; j++ {
			scale := n.rbf(x, y, j, i)
			sum += scale * in[j+i*n.Width]
		}
	}
	return sum
}

func (n *rbfMatrix) rbf(x1, y1, x2, y2 int) float64 {
	if n.rbfCache == nil {
		maxDist := int(-2*math.Log(2*n.Variance*epsilon) + 4)
		n.rbfCache = make([][]float64, maxDist)
		for i := range n.rbfCache {
			n.rbfCache[i] = make([]float64, maxDist)
			for j := range n.rbfCache[i] {
				dist := float64(i*i + j*j)
				rbf := math.Exp(-dist / (2 * n.Variance))
				n.rbfCache[i][j] = rbf
			}
		}
	}
	return n.rbfCache[absInt(x1-x2)][absInt(y1-y2)]
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func pointBounds(x, y, variance float64) (minX, minY, maxX, maxY int) {
	maxDist := -math.Log(2 * variance * epsilon)
	minX = int(math.Max(0, x-maxDist))
	minY = int(math.Max(0, y-maxDist))
	maxX = int(x + maxDist + 1)
	maxY = int(y + maxDist + 1)
	return
}
