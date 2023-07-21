package timescale

import "testing"

func TestIntTime2Time(t *testing.T) {
	//timeInt := 93000
	timeInt := 93000111
	t.Logf("%s", IntTime2Time(timeInt))
}
