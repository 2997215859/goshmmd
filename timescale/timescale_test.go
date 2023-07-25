package timescale

import (
	"fmt"
	"testing"
)

func TestNewTimeScale(t *testing.T) {
	timescale := NewTimeScale("09:30:00", "15:00:00", 60)
	fmt.Println(timescale)
}
