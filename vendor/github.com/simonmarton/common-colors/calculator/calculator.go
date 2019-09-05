package calculator

import (
	"fmt"
	"math"

	"github.com/simonmarton/common-colors/color"
	"github.com/simonmarton/common-colors/models"
)

const defaultTransparencyTreshold uint8 = 10
const defaultIterationCount int8 = 8
const defaultMinLuminance float64 = 0
const defaultMaxLuminance float64 = 1
const defaultDistanceThreshold float64 = 50
const defaultMinSaturation float64 = .3
const defaultAlgorithm string = "simple"

// Calculator can group common colors
type Calculator struct {
	config models.CalculatorConfig
}

// New Calculator instance
func New(c models.CalculatorConfig) *Calculator {
	if c.TransparencyTreshold <= 0 {
		c.TransparencyTreshold = defaultTransparencyTreshold
	}

	if c.IterationCount <= 0 {
		c.IterationCount = defaultIterationCount
	}

	if c.MinLuminance < 0 {
		c.MinLuminance = defaultMinLuminance
	}

	if c.MaxLuminance > 1 || c.MaxLuminance < .0001 {
		c.MaxLuminance = defaultMaxLuminance
	}

	// TODO write test
	if c.MinSaturation < 0 || c.MinSaturation > 1 {
		c.MinSaturation = defaultMinSaturation
	}

	// TODO write test
	// 441.67 is the max distance in a cube with 255 long sides
	if c.DistanceThreshold < 1 || c.DistanceThreshold > 441.67 {
		c.DistanceThreshold = defaultDistanceThreshold
	}

	if c.Algorithm == "" {
		c.Algorithm = defaultAlgorithm
	}

	return &Calculator{config: c}
}

// GetCommonColors ...
func (c Calculator) GetCommonColors(colors []color.Color) ([]color.Color, [][]color.Color) {
	fmt.Printf("Colors length: %d\n", len(colors))
	stepsOfColors := [][]color.Color{colors}

	colors = c.removeInvalidColors(colors)
	fmt.Printf("Colors length after: %d\n", len(colors))

	stepsOfColors = append(stepsOfColors, colors)

	for i := int8(0); i < c.config.IterationCount; i++ {
		threshold := c.config.DistanceThreshold*float64(i)/float64(c.config.IterationCount-1) + 10
		colors = c.groupByThreshold(colors, threshold)

		stepsOfColors = append(stepsOfColors, color.Sort(colors))
	}

	color.Sort(colors)

	return colors, stepsOfColors
}

// GenrateGradientColors ...
func (c Calculator) GenrateGradientColors(colors []color.Color) (result []string) {
	totalWeight := 0

	for _, col := range colors {
		totalWeight += col.Weight
	}

	mainColor := colors[0]
	var secondaryColor string

	if len(colors) >= 2 {
		mp := float64(mainColor.Weight) / float64(totalWeight)
		for _, col := range colors[1:] {
			p := float64(col.Weight) / float64(totalWeight)

			// Too big weight diff compared to the main color, ignore
			// if mp-p > .4 {
			if mp/p > 2.5 {
				break
			}

			// Colors have similar hue, diff is less than 45deg
			hd := math.Abs(mainColor.Hue() - col.Hue())

			if hd < .125 {
				secondaryColor = col.ToHex()
				break
			}
		}
	}

	if secondaryColor == "" {
		// Calculate from main color
		h, s, l, _ := mainColor.ToHSLA()

		// Rotate 10deg
		h += .025 // TODO config
		s = math.Max(c.config.MinSaturation, s-.2)

		secondaryColor = color.NewFromHSL(h, s, l).ToHex()
	}

	return []string{mainColor.ToHex(), secondaryColor}
}

func (c Calculator) removeInvalidColors(colors []color.Color) (result []color.Color) {
	for _, col := range colors {
		l := col.Luminance()
		s := col.Saturation()

		// todo
		if col.A <= c.config.TransparencyTreshold ||
			l < c.config.MinLuminance ||
			l > c.config.MaxLuminance ||
			s < c.config.MinSaturation {
			continue
		}

		result = append(result, col)
	}

	return result
}

func (c Calculator) groupByThreshold(colors []color.Color, threshold float64) (result []color.Color) {
	for len(colors) > 1 {
		sample := colors[0]

		similarColors := []color.Color{sample}
		remainingColors := []color.Color{}

		for _, color := range colors[1:] {
			var d float64
			switch c.config.Algorithm {
			case "simple":
				d = sample.Distance(color)
			case "yiq":
				d = sample.YIQDistance(color)
			default:
				panic("Not supported Algorithm")
			}

			if d < threshold {
				similarColors = append(similarColors, color)
			} else {
				remainingColors = append(remainingColors, color)
			}
		}

		colors = remainingColors

		result = append(result, color.Average(similarColors))
	}

	// TODO? colors at this point can still have some very different colors
	// should we append them too?
	// meh, why not?
	return append(result, colors...)
}

// Nope
func (c Calculator) groupByNearestNeighbor(colors []color.Color) (result []color.Color) {
	for len(colors) > 1 {
		sample := colors[0]
		distance := math.MaxFloat64
		closestColorIdx := -1
		for idx, color := range colors[1:] {
			d := sample.Distance(color)
			if d < distance {
				distance = d
				closestColorIdx = idx + 1
			}
		}

		// fmt.Printf("groupByNearestNeighbor closestColorIdx: %d\n", closestColorIdx)

		result = append(result, sample.Average(colors[closestColorIdx]))
		// Remove remove sample (first item) and the closes color
		colors = append(colors[1:closestColorIdx], colors[closestColorIdx+1:]...)
	}

	if len(colors) > 0 {
		result = append(result, colors...)
	}

	return result
}
