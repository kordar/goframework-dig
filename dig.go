package goframework_dig

import (
	"go.uber.org/dig"
)

type DigIns struct {
	name string
	ins  *dig.Container
}

func NewDigIns(name string, options ...dig.Option) *DigIns {
	client := dig.New(options...)
	return &DigIns{name: name, ins: client}
}

func (c DigIns) GetName() string {
	return c.name
}

func (c DigIns) GetInstance() interface{} {
	return c.ins
}

func (c DigIns) Close() error {
	return nil
}
