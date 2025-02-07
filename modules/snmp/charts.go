package snmp

import (
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

//TODO: probably not needed
var defaultSNMPchart = module.Chart{
	ID:       "snmp_%d",
	Title:    "%s",
	Units:    "kilobits/s",
	Type:     "area",
	Priority: 1,
	Fam:      "ports",
}

// newChart creates news chart based on 'ChartsConfig', 'id' and 'oidIndex'
// parameters. oidIndex is optional param, which decided whether to add an
// index to OID value or not.
func newChart(id int, oidIndex *int, s ChartsConfig) (*module.Chart, error) {
	c := defaultSNMPchart.Copy()
	c.ID = fmt.Sprintf(c.ID, id)
	c.Title = fmt.Sprintf(c.Title, s.Title)

	if oidIndex != nil {
		c.ID = fmt.Sprintf("%s_%d", c.ID, *oidIndex)
		c.Title = fmt.Sprintf("%s %d", c.Title, *oidIndex)
	}

	if s.Family != nil {
		c.Fam = *s.Family
	}

	if s.Units != nil {
		c.Units = *s.Units
	}

	if s.Type != nil {
		c.Type = module.ChartType(*s.Type)
	}

	c.Priority = s.Priority
	for _, d := range s.Dimensions {
		oid := d.OID
		if oidIndex != nil {
			oid = fmt.Sprintf("%s.%d", oid, *oidIndex)
		}
		dim := &module.Dim{
			Name: d.Name,
			ID:   oid,
		}
		if d.Algorithm != nil {
			dim.Algo = module.DimAlgo(*d.Algorithm)
		}
		if d.Multiplier != nil {
			dim.Mul = *d.Multiplier
		}
		if d.Divisor != nil {
			dim.Div = *d.Divisor
		}

		if err := c.AddDim(dim); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func newCharts(configs []ChartsConfig) (*module.Charts, error) {
	charts := &module.Charts{}
	for i, cfg := range configs {
		if cfg.MultiplyRange != nil {
			for j := cfg.MultiplyRange[0]; j <= cfg.MultiplyRange[1]; j++ {
				chart, err := newChart(i, &j, cfg)
				if err != nil {
					return nil, err
				}
				if err = charts.Add(chart); err != nil {
					return nil, err
				}
			}
		} else {
			chart, err := newChart(i, nil, cfg)
			if err != nil {
				return nil, err
			}
			if err = charts.Add(chart); err != nil {
				return nil, err
			}
		}
	}
	return charts, nil
}
