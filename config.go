package main

import (
	"github.com/hashicorp/packer/template/interpolate"
)

type FlasherConfig struct {
	Device      string `mapstructure:"device"`
	Interactive bool   `mapstructure:"interactive"`
	BlockSize   int    `mapstructure:"block_size"`
}

func (c *FlasherConfig) Prepare(ctx *interpolate.Context) (warnings []string, errs []error) {
	return warnings, errs
}
