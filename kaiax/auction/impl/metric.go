package impl

import "github.com/rcrowley/go-metrics"

var (
	numBidsGauge         = metrics.NewRegisteredGauge("kaiax/auction/bidpool/num/bids", nil)
	numBidRequestCounter = metrics.NewRegisteredCounter("kaiax/auction/bidpool/num/bidreqs", nil)
)
