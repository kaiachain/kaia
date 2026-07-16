// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package setup

import (
	"bytes"
	"fmt"
	"text/template"
)

type PrometheusConfig struct {
	CnIps []string
	PnIps []string
}

func NewPrometheusConfig(cnNum int, cnNetworkIp string) PrometheusConfig {
	var cnIps []string
	for i := 1; i <= cnNum; i++ {
		cnIps = append(cnIps, fmt.Sprintf("%s.%d", cnNetworkIp, 100+i))
	}

	return PrometheusConfig{
		CnIps: cnIps,
	}
}

func (pConfig PrometheusConfig) String() string {
	tmpl, err := template.New("PrometheusConfig").Parse(pTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template, %v", err)
		return ""
	}

	result := new(bytes.Buffer)
	err = tmpl.Execute(result, pConfig)
	if err != nil {
		fmt.Printf("Failed to render template, %v", err)
		return ""
	}

	return result.String()
}

var pTemplate = `global:
  scrape_interval:     15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'klaytn'
    static_configs:
    - targets:
      {{- range .CnIps }}
      - "{{ . }}:61001"
      {{- end }}
      labels:
        job: 'klaytn-cn'
    {{- if gt (len .PnIps) 0 }}
    - targets:
      {{- range .PnIps }}
      - "{{ . }}:61001"
      {{- end }}
      labels:
        job: 'klaytn-pn'
    {{ end }}
`
