package rules

import (
	"math"
	"time"
)

type SeasonalBaseline struct {
	Hour    [24]float64
	Weekday [7]float64
	Count   [31]int
}

func (b *SeasonalBaseline) Add(at time.Time, value float64) {
	b.Hour[at.Hour()] += value
	b.Weekday[int(at.Weekday())] += value
	b.Count[at.Day()]++
}
func (b SeasonalBaseline) HourAverage(hour int) float64 {
	if hour < 0 || hour >= 24 {
		return 0
	}
	return b.Hour[hour]
}
func (b SeasonalBaseline) WeekdayAverage(day time.Weekday) float64 { return b.Weekday[int(day)] }
func SeasonalDistance(a, b time.Time) float64 {
	hours := math.Abs(a.Sub(b).Hours())
	if hours > 12 {
		hours = 24 - hours
	}
	return hours
}
func IsNight(at time.Time) bool { return at.Hour() < 6 || at.Hour() >= 22 }
