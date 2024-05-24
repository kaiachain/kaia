// Copyright 2024 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.
package rpc

import "github.com/rcrowley/go-metrics"

var (
	rpcTotalRequestsCounter    = metrics.NewRegisteredCounter("rpc/counts/total", nil)
	rpcSuccessResponsesCounter = metrics.NewRegisteredCounter("rpc/counts/success", nil)
	rpcErrorResponsesCounter   = metrics.NewRegisteredCounter("rpc/counts/errors", nil)
	rpcPendingRequestsCount    = metrics.NewRegisteredCounter("rpc/counts/pending", nil)

	wsSubscriptionReqCounter   = metrics.NewRegisteredCounter("ws/counts/subscription/request", nil)
	wsUnsubscriptionReqCounter = metrics.NewRegisteredCounter("ws/counts/unsubscription/request", nil)
	wsConnCounter              = metrics.NewRegisteredCounter("ws/counts/connections/total", nil)
)
