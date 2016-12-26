package rbfscale

import (
	"math"

	"github.com/unixpickle/num-analysis/linalg"
)

const epsilon = 1e-8

type rbfCacheEntry struct {
	InIndex int
	Coeff   float64
}

// A rbfMatrix acts on behalf of an RBF layer with centers
// arranged on the points in an image.
// The Variance is used to divide the argument to exp,
// giving the formula exp(||x-c||/(2*Variance)).
type rbfMatrix struct {
	Variance float64
	Width    int
	Height   int

	cache [][]rbfCacheEntry
}

// Dim returns the total number of pixels.
func (n *rbfMatrix) Dim() int {
	return n.Width * n.Height
}

// Apply performs matrix multiplication by the positive
// definite matrix Gij = Rho(Pi-Pj), where Pi is the i-th
// point in the image.
func (n *rbfMatrix) Apply(in linalg.Vector) linalg.Vector {
	if n.cache == nil {
		n.buildCache()
	}
	result := make(linalg.Vector, n.Dim())
	for i, caches := range n.cache {
		for _, c := range caches {
			result[i] += in[c.InIndex] * c.Coeff
		}
	}
	return result
}

func (n *rbfMatrix) buildCache() {
	n.cache = make([][]rbfCacheEntry, 0, n.Width*n.Height)
	for y := 0; y < n.Height; y++ {
		for x := 0; x < n.Width; x++ {
			n.cache = append(n.cache, n.buildPixelCache(x, y))
		}
	}
}

func (n *rbfMatrix) buildPixelCache(x, y int) []rbfCacheEntry {
	minX, minY, maxX, maxY := pointBounds(float64(x), float64(y), n.Variance)
	var res []rbfCacheEntry
	for i := minY; i < maxY && i < n.Height; i++ {
		for j := minX; j < maxX && j < n.Width; j++ {
			scale := n.rbf(x, y, j, i)
			res = append(res, rbfCacheEntry{
				InIndex: j + i*n.Width,
				Coeff:   scale,
			})
		}
	}
	return res
}

func (n *rbfMatrix) rbf(x1, y1, x2, y2 int) float64 {
	dist := float64((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
	return math.Exp(-dist / (2 * n.Variance))
}

func pointBounds(x, y, variance float64) (minX, minY, maxX, maxY int) {
	maxDist := math.Sqrt(-math.Log(epsilon)*2*variance) + 1
	minX = int(math.Max(0, x-maxDist))
	minY = int(math.Max(0, y-maxDist))
	maxX = int(x + maxDist + 1)
	maxY = int(y + maxDist + 1)
	return
}
