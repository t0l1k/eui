package eui

import (
	"fmt"
)

type UiConfig struct {
	pref   *Preferences
	locale *Locale
	theme  *Theme
}

func NewConfig() *UiConfig {
	c := &UiConfig{}
	return c
}

func (c *UiConfig) SetPref(value *Preferences) {
	c.pref = value
}

func (c *UiConfig) SetLocale(value *Locale) {
	c.locale = value
}

func (c *UiConfig) SetTheme(value *Theme) {
	c.theme = value
}

type Preferences map[string]interface{}

func NewPreferences() Preferences {
	return make(Preferences)
}

func (p Preferences) Get(set string) (value interface{}) {
	return p[set]
}

func (p Preferences) Set(set string, value interface{}) {
	p[set] = value
}

func (p Preferences) String() string {
	s := ""
	for k, v := range p {
		s += fmt.Sprintf("%v: %v\n", k, v)
	}
	return s
}

type Locale map[string]string

func NewLocale() Locale {
	return make(Locale)
}

func (l Locale) Get(value string) string {
	return l[value]
}

func (l Locale) Set(set, value string) {
	l[set] = value
}

func (l Locale) String() string {
	s := ""
	for k, v := range l {
		s += fmt.Sprintf("%v: %v\n", k, v)
	}
	return s
}
