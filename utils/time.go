package utils

import (
	"fmt"
	"math"
)

func DurationToString(duration float64) string {
	if duration <= 20.0 {
		return "<0:20"
	}
	minutes := math.Floor(duration / 60.0)
	second := duration - minutes*60.0
	return fmt.Sprintf("%.0f:%02.0f", minutes, second)
}
