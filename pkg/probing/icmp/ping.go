package icmp

import (
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type Statistics struct {
	probing.Statistics
}

func Ping(address string, count int, duration time.Duration) *Statistics {
	pinger, err := probing.NewPinger(address)
	if err != nil {
		return nil
	}
	pinger.Count = count
	pinger.Timeout = duration
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return nil
	}
	stat := pinger.Statistics() // get send/receive/duplicate/rtt stats
	return &Statistics{
		*stat,
	}
}
