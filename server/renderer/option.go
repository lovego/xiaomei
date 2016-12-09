package renderer

import (
	"html/template"
)

type LayoutDataGetter interface {
	GetLayoutData(layout string, data interface{}) interface{}
}

type O struct {
	NoLayout         bool
	Layout           string
	LayoutDataGetter LayoutDataGetter
	Funcs            template.FuncMap
}

func (option O) Process(renderer *Renderer) O {
	if option.NoLayout {
		return option
	}
	if option.Layout == `` {
		option.Layout = renderer.Layout
	}
	return option
}

func (option O) HasLayout() bool {
	return !option.NoLayout && option.Layout != ``
}

func (option O) LayoutData(data interface{}) interface{} {
	if option.LayoutDataGetter == nil {
		return data
	}
	return option.LayoutDataGetter.GetLayoutData(option.Layout, data)
}
