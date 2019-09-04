package models

// CalculatorConfig defines the parameters for the calculator to use
type CalculatorConfig struct {
	TransparencyTreshold uint8   `json:"transparencyTreshold"`
	IterationCount       int8    `json:"iterationCount"`
	MinLuminance         float64 `json:"minLuminance"`
	MaxLuminance         float64 `json:"maxLuminance"`
	DistanceThreshold    float64 `json:"distanceThreshold"`
	MinSaturation        float64 `json:"minSaturation"`
	Algorithm            string  `json:"algorithm"`
}
