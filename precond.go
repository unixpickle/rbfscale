package rbfscale

import (
	"math"

	"github.com/unixpickle/num-analysis/linalg"
	"github.com/unixpickle/num-analysis/linalg/svd"
)

const precondRadius = 3

// A preconditioner preconditions conjugate gradients for
// solving the interpolation system.
type preconditioner struct {
	row linalg.Vector
	dim int
}

func newPreconditioner(width, height int, variance float64) *preconditioner {
	n := precondRadius*2 + 1
	mat := linalg.NewMatrix(n, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			mag := math.Exp(-math.Pow(float64(i-j), 2) / (2 * variance))
			mat.Set(i, j, mag)
		}
	}
	v, d, u := svd.Decompose(mat)
	for i := 0; i < n; i++ {
		d.Set(i, i, 1/math.Sqrt(float64(d.Get(i, i))))
	}
	inv := v.Mul(d).Mul(u)
	return &preconditioner{
		row: inv.Data[precondRadius*inv.Cols : (precondRadius+1)*inv.Cols],
		dim: width * height,
	}
}

func (p *preconditioner) Dim() int {
	return p.dim
}

func (p *preconditioner) Apply(in linalg.Vector) linalg.Vector {
	res := make(linalg.Vector, len(in))
	for i := range in {
		offset := (len(p.row) - 1) / 2
		for j := i - offset; j <= i+offset && j < p.dim; j++ {
			if j < 0 {
				continue
			}
			coeff := p.row[j-i+offset]
			res[i] += in[j] * coeff
		}
	}
	return res
}
