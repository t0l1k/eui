package eui

import (
	"fmt"
)

func DefaultSettings() *Setting {
	s := NewSetting()
	s.Set(UiFullscreen, false)
	return &s
}

type SettingName int

const (
	UiFullscreen SettingName = iota
)

type Setting map[SettingName]interface{}

func NewSetting() Setting {
	return make(Setting)
}

func (p Setting) Get(set SettingName) (value interface{}) {
	return p[set]
}

func (p Setting) Set(set SettingName, value interface{}) {
	p[set] = value
}

func (p Setting) String() string {
	s := ""
	for k, v := range p {
		s += fmt.Sprintf("%v: %v\n", k, v)
	}
	return s
}
