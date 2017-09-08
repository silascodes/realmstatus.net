package main

import (
	"strings"
	"common"
)

type Region struct {
	common.Region
}

func (this Region) GetValue() string {
	return strings.ToLower(string(this.Region))
}

func (this Region) GetDisplay() string {
	return strings.ToUpper(string(this.Region))
}