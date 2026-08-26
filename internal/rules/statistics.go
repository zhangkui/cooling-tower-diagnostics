package rules

import "math"

type Statistics struct {
	Count    int
	Sum      float64
	Mean     float64
	Variance float64
	StdDev   float64
	Min      float64
	Max      float64
	P50      float64
	P95      float64
}

func Describe(values []float64) Statistics {
	stats := Statistics{Count: len(values)}
	if len(values) == 0 {
		return stats
	}
	stats.Min, stats.Max = values[0], values[0]
	for _, value := range values {
		stats.Sum += value
		if value < stats.Min {
			stats.Min = value
		}
		if value > stats.Max {
			stats.Max = value
		}
	}
	stats.Mean = stats.Sum / float64(stats.Count)
	for _, value := range values {
		delta := value - stats.Mean
		stats.Variance += delta * delta
	}
	stats.Variance /= float64(stats.Count)
	stats.StdDev = math.Sqrt(stats.Variance)
	stats.P50 = Percentile(values, 0.50)
	stats.P95 = Percentile(values, 0.95)
	return stats
}

func ZScore(value float64, stats Statistics) float64 {
	if stats.StdDev == 0 {
		return 0
	}
	return (value - stats.Mean) / stats.StdDev
}
func IsOutlier(value float64, stats Statistics, limit float64) bool {
	return math.Abs(ZScore(value, stats)) >= limit
}
func CoefficientOfVariation(stats Statistics) float64 {
	if stats.Mean == 0 {
		return 0
	}
	return stats.StdDev / math.Abs(stats.Mean)
}
