package rbfscale

import (
	"image"
	"image/color"
	"math"

	"github.com/unixpickle/num-analysis/conjgrad"
	"github.com/unixpickle/num-analysis/linalg"
)

const (
	weightEpsilon = 1e-2
	numChannels   = 3
)

// An Interp uses an RBF network to interpolate between
// the pixels of an image.
type Interp struct {
	variance float64
	weights  [numChannels]linalg.Vector
	width    int
	height   int
}

// NewInterp learns an RBF network from the image, using
// the variance as a parameter to the internal RBFs.
func NewInterp(img image.Image, variance float64) *Interp {
	mat := &rbfMatrix{
		Variance: variance,
		Width:    img.Bounds().Dx(),
		Height:   img.Bounds().Dy(),
	}
	res := &Interp{
		variance: mat.Variance,
		width:    mat.Width,
		height:   mat.Height,
	}
	for channel := 0; channel < numChannels; channel++ {
		colors := getColorChannel(img, channel)
		res.weights[channel] = conjgrad.SolvePrec(mat, colors, weightEpsilon)
	}
	return res
}

// Width returns the width of the original image.
func (n *Interp) Width() int {
	return n.width
}

// Height returns the height of the original image.
func (n *Interp) Height() int {
	return n.height
}

// At performs color interpolation to compute the color at
// the given coordinate, starting at (0, 0).
func (n *Interp) At(x, y float64) color.Color {
	channelSum := [numChannels]float64{}
	minX, minY, maxX, maxY := pointBounds(x, y, n.variance)
	for i := minY; i < maxY && i < n.height; i++ {
		for j := minX; j < maxX && j < n.width; j++ {
			dist := (float64(j)-x)*(float64(j)-x) + (float64(i)-y)*(float64(i)-y)
			rbf := math.Exp(-dist / (2 * n.variance))
			for ch := range channelSum[:] {
				channelSum[ch] += rbf * n.weights[ch][j+i*n.width]
			}
		}
	}
	for i, x := range channelSum[:] {
		channelSum[i] = math.Max(0, math.Min(1, x))
	}
	return color.RGBA{
		R: uint8(channelSum[0] * 0xff),
		G: uint8(channelSum[1] * 0xff),
		B: uint8(channelSum[2] * 0xff),
		A: 0xff,
	}
}

// Image produces an image of the given dimensions via
// interpolation.
func (n *Interp) Image(width, height int) image.Image {
	xScale := float64(n.width) / float64(width)
	yScale := float64(n.height) / float64(height)
	res := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		sourceY := float64(y) * yScale
		for x := 0; x < width; x++ {
			sourceX := float64(x) * xScale
			res.Set(x, y, n.At(sourceX, sourceY))
		}
	}
	return res
}

func getColorChannel(img image.Image, idx int) linalg.Vector {
	var res linalg.Vector
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r, g, b, _ := color.RGBAModel.Convert(img.At(x, y)).RGBA()
			nums := []uint32{r, g, b}
			res = append(res, float64(nums[idx])/0xffff)
		}
	}
	return res
}
