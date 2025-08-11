// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from node/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package node

import (
	"fmt"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"
)

// DebugNodeAPI is the collection of debugging related API methods exposed over
// both secure and unsecure RPC channels.
type DebugNodeAPI struct {
	node *Node // Node interfaced by this API
}

// NewDebugNodeAPI creates a new API definition for the public debug methods
// of the node itself.
func NewDebugNodeAPI(node *Node) *DebugNodeAPI {
	return &DebugNodeAPI{node: node}
}

// Metrics retrieves all the known system metric collected by the node.
func (api *DebugNodeAPI) Metrics(raw bool) (map[string]interface{}, error) {
	// Create a rate formatter
	units := []string{"", "K", "M", "G", "T", "E", "P"}
	round := func(value float64, prec int) string {
		unit := 0
		for value >= 1000 {
			unit, value, prec = unit+1, value/1000, 2
		}
		return fmt.Sprintf(fmt.Sprintf("%%.%df%s", prec, units[unit]), value)
	}
	format := func(total float64, rate float64) string {
		return fmt.Sprintf("%s (%s/s)", round(total, 0), round(rate, 2))
	}
	// Iterate over all the metrics, and just dump for now
	counters := make(map[string]interface{})
	metrics.DefaultRegistry.Each(func(name string, metric interface{}) {
		// Create or retrieve the counter hierarchy for this metric
		root, parts := counters, strings.Split(name, "/")
		for _, part := range parts[:len(parts)-1] {
			if _, ok := root[part]; !ok {
				root[part] = make(map[string]interface{})
			}
			root = root[part].(map[string]interface{})
		}
		name = parts[len(parts)-1]

		// Fill the counter with the metric details, formatting if requested
		if raw {
			switch metric := metric.(type) {
			case metrics.Counter:
				root[name] = map[string]interface{}{
					"Overall": float64(metric.Count()),
				}

			case metrics.Meter:
				root[name] = map[string]interface{}{
					"AvgRate01Min": metric.Rate1(),
					"AvgRate05Min": metric.Rate5(),
					"AvgRate15Min": metric.Rate15(),
					"MeanRate":     metric.RateMean(),
					"Overall":      float64(metric.Count()),
				}

			case metrics.Timer:
				root[name] = map[string]interface{}{
					"AvgRate01Min": metric.Rate1(),
					"AvgRate05Min": metric.Rate5(),
					"AvgRate15Min": metric.Rate15(),
					"MeanRate":     metric.RateMean(),
					"Mean":         metric.Mean(),
					"StdDev":       metric.StdDev(),
					"Overall":      float64(metric.Count()),
					"Percentiles": map[string]interface{}{
						"5":  metric.Percentile(0.05),
						"20": metric.Percentile(0.2),
						"50": metric.Percentile(0.5),
						"80": metric.Percentile(0.8),
						"95": metric.Percentile(0.95),
					},
				}

			default:
				root[name] = "Unknown metric type"
			}
		} else {
			switch metric := metric.(type) {
			case metrics.Counter:
				root[name] = map[string]interface{}{
					"Overall": float64(metric.Count()),
				}

			case metrics.Meter:
				root[name] = map[string]interface{}{
					"Avg01Min": format(metric.Rate1()*60, metric.Rate1()),
					"Avg05Min": format(metric.Rate5()*300, metric.Rate5()),
					"Avg15Min": format(metric.Rate15()*900, metric.Rate15()),
					"Overall":  format(float64(metric.Count()), metric.RateMean()),
				}

			case metrics.Timer:
				root[name] = map[string]interface{}{
					"Avg01Min": format(metric.Rate1()*60, metric.Rate1()),
					"Avg05Min": format(metric.Rate5()*300, metric.Rate5()),
					"Avg15Min": format(metric.Rate15()*900, metric.Rate15()),
					"Overall":  format(float64(metric.Count()), metric.RateMean()),
					"Maximum":  time.Duration(metric.Max()).String(),
					"Minimum":  time.Duration(metric.Min()).String(),
					"Percentiles": map[string]interface{}{
						"5":  time.Duration(metric.Percentile(0.05)).String(),
						"20": time.Duration(metric.Percentile(0.2)).String(),
						"50": time.Duration(metric.Percentile(0.5)).String(),
						"80": time.Duration(metric.Percentile(0.8)).String(),
						"95": time.Duration(metric.Percentile(0.95)).String(),
					},
				}

			default:
				root[name] = "Unknown metric type"
			}
		}
	})
	return counters, nil
}
