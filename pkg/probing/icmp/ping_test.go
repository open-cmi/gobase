package icmp

import (
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	s := Ping("114.114.114.114", 3, 5*time.Second)
	if s == nil {
		t.Errorf("ping error")
		return
	}
	t.Logf("avg rtt is %d\n", s.AvgRtt)
}
