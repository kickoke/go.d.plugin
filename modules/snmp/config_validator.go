package snmp

import (
	"errors"
	"fmt"

	"github.com/netdata/go.d.plugin/agent/module"
)

func appendError(err error, msg string) error {
	var e error
	if err == nil {
		if msg == "" {
			return nil
		}
		e = errors.New(msg)
	} else {
		if msg == "" {
			return err
		}
		e = fmt.Errorf("%w %s", err, msg)
	}
	return e
}

func (d Dimension) validateConfig(index_chart, index int) error {
	var err error
	err = nil
	if d.Name == "" {
		err = appendError(err, fmt.Sprintf("invalid or missing value: charts[%d].dimension[%d].name;", index_chart, index))
	}
	if d.OID == "" {
		err = appendError(err, fmt.Sprintf("missing value: charts[%d].dimension[%d].oid;", index_chart, index))
	}
	if d.Algorithm != nil {
		if *d.Algorithm == "" ||
			(*d.Algorithm != string(module.Incremental) &&
				*d.Algorithm != string(module.PercentOfIncremental) &&
				*d.Algorithm != string(module.PercentOfAbsolute)) {
			err = appendError(err, fmt.Sprintf("invalid or missing value: charts[%d].dimension[%d].algorithm;", index_chart, index))
		}
	}
	if d.Multiplier != nil {
		if *d.Multiplier == 0 {
			err = appendError(err, fmt.Sprintf("integer set to 0: charts[%d].dimension[%d].multiplier;", index_chart, index))
		}
	}
	if d.Divisor != nil {
		if *d.Divisor == 0 {
			err = appendError(err, fmt.Sprintf("integer set to 0: charts[%d].dimension[%d].divisor;", index_chart, index))
		}
	}
	return err
}

func (u User) validateConfig() error {
	var err error
	err = nil
	if u.Name == "" {
		err = appendError(err, "missing value: user.name;")
	}
	if u.Level < 1 || u.Level > 3 {
		err = appendError(err, fmt.Sprintf("invalid range of value(%d): user.level;", u.Level))
	}
	if u.PrivProto < 1 || u.PrivProto > 2 {
		err = appendError(err, fmt.Sprintf("invalid range of value(%d): user.priv_proto;", u.PrivProto))
	}
	if u.AuthProto < 1 || u.AuthProto > 3 {
		err = appendError(err, fmt.Sprintf("invalid range of value(%d): user.auth_proto;", u.AuthProto))
	}
	return err
}

func (c ChartsConfig) validateConfig(index_chart int) error {
	var err error
	err = nil
	if c.Title == "" {
		err = appendError(err, fmt.Sprintf("missing value: charts[%d].title;", index_chart))
	}
	if c.Dimensions == nil {
		err = appendError(err, fmt.Sprintf("missing value: charts[%d].dimensions;", index_chart))
	} else {
		for i, d := range c.Dimensions {
			if e := d.validateConfig(index_chart, i); e != nil {
				err = appendError(err, e.Error())
			}
		}
	}

	if c.MultiplyRange != nil {
		if len(c.MultiplyRange) != 2 {
			err = appendError(err, fmt.Sprintf("invalid range: charts[%d].multiply_range;", index_chart))
		} else {
			if c.MultiplyRange[0] >= c.MultiplyRange[1] || c.MultiplyRange[0] < 0 {
				err = appendError(err, fmt.Sprintf("invalid range: charts[%d].multiply_range;", index_chart))
			}
		}
	}
	return err
}

func (o Options) validateConfig() error {
	var err error
	err = nil
	if o.Port <= 0 && o.Port > 65535 {
		err = appendError(err, fmt.Sprintf("invalid range of value(%d): options.port;", o.Port))
	}
	if o.Version < 1 || o.Version > 3 {
		err = appendError(err, fmt.Sprintf("invalid range of value(%d): options.versions;", o.Version))
	}
	if o.Retries < 1 || o.Retries > 100 {
		err = appendError(err, fmt.Sprintf("invalid range of value(%d): options.retries;", o.Retries))
	}
	if o.Timeout < 1 {
		err = appendError(err, fmt.Sprintf("invalid value(%d): options.timeout;", o.Timeout))
	}
	if o.MaxOIDs <= 0 {
		err = appendError(err, fmt.Sprintf("invalid value(%d): max_request_size;", o.MaxOIDs))
	}

	return err
}

func (s SNMP) validateConfig() error {
	var err error
	err = nil
	if s.ChartInput == nil {
		err = appendError(err, "charts missing from config")
	} else {
		for i, chartIn := range s.ChartInput {
			if e := chartIn.validateConfig(i); e != nil {
				err = appendError(err, e.Error())
			}
		}
	}

	if s.Options != nil {
		if e := s.Options.validateConfig(); e != nil {
			err = appendError(err, e.Error())
		}
		if s.Options.Version == 3 {
			if s.User == nil {
				err = appendError(err, "SNMP v3 missing user credentials;")
			}
		} else {
			if s.Community == nil {
				err = appendError(err, "SNMP v1/2 missing community value;")
			}
		}
	}

	if s.User != nil {
		if e := s.User.validateConfig(); e != nil {
			err = appendError(err, e.Error())
		}
	}

	if s.Config.UpdateEvery <= 0 {
		err = appendError(err, fmt.Sprintf("invalid value(%d): update_every;", s.Config.UpdateEvery))
	}

	if s.Config.Name == "" {
		err = appendError(err, "missing value: name;")
	}

	return err
}
