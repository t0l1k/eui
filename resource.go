package eui

import "log"

const (
	FontDefault string = "default"
)

type ResourceManager struct {
	fonts map[string]*Font
}

func NewResourceManager() *ResourceManager { return &ResourceManager{fonts: make(map[string]*Font)} }
func (r *ResourceManager) LoadFont(name string, data []byte, size int) *ResourceManager {
	var err error
	r.fonts[name], err = NewFont(name, data, size)
	if err != nil {
		panic(err)
	}
	log.Println("ResourceManager:LoadFont:", name, size)
	return r
}
func (a *ResourceManager) FontDefault() *Font { return a.fonts[FontDefault] }
