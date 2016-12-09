package renderer

import (
	// "fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"
	"strings"
)

type Tmpl struct {
	*template.Template
	loaded map[string]bool
}

type Renderer struct {
	Root   string
	Layout string
	Tmpls  map[string]*Tmpl
	Funcs  template.FuncMap
}

func New(root, layout string, cache bool, funcs template.FuncMap) *Renderer {
	var tmpls map[string]*Tmpl
	if cache {
		tmpls = make(map[string]*Tmpl)
	}
	return &Renderer{
		Root:   path.Clean(root),
		Layout: path.Clean(layout),
		Tmpls:  tmpls,
		Funcs:  funcs,
	}
}

func (r *Renderer) Render(wr io.Writer, name string, data interface{}, option O) {
	option = option.Process(r)
	tmpl := r.getTemplate(name, option)
	var err error
	if option.HasLayout() {
		err = tmpl.Template.ExecuteTemplate(wr, option.Layout, option.LayoutData(data))
	} else {
		err = tmpl.Template.Execute(wr, data)
	}
	if err != nil {
		panic(err)
	}
}

func (r *Renderer) getTemplate(name string, option O) *Tmpl {
	name = cleanName(name)
	tmpl := r.Tmpls[name]
	if tmpl == nil {
		tmpl = &Tmpl{template.New(``), make(map[string]bool)}
		if r.Funcs != nil {
			tmpl.Funcs(r.Funcs)
		}
		if option.Funcs != nil {
			tmpl.Funcs(option.Funcs)
		}
		parseTemplate(tmpl.Template, name, r.Root, r.Root, tmpl.loaded)
		if r.Tmpls != nil {
			r.Tmpls[name] = tmpl
		}
	}
	if option.HasLayout() {
		layout := option.Layout
		if tmpl.Lookup(layout) == nil {
			parseTemplate(tmpl.Template.New(layout), layout, r.Root, r.Root, tmpl.loaded)
		}
	}
	return tmpl
}

// 关联模板中不能含有同名模板
func parseTemplate(templ *template.Template, name, base, root string, loaded map[string]bool) {
	var p string
	if path.IsAbs(name) {
		p = path.Join(root, name)
	} else {
		p = path.Join(base, name)
	}
	content, err := ioutil.ReadFile(p + `.tmpl`)
	if err != nil {
		panic(err)
	}
	template.Must(templ.Parse(string(content)))

	var new []string
	for _, t := range templ.Templates() { // 所有已关联的模板，包括自己
		nam := t.Name()
		if nam != name && strings.IndexByte(nam, '/') >= 0 && !loaded[nam] {
			loaded[nam] = true
			new = append(new, nam)
		}
	}

	base = path.Dir(p)
	for _, nam := range new {
		parseTemplate(templ.New(nam), nam, base, root, loaded)
	}
}

func cleanName(name string) string {
	if strings.IndexByte(name, '/') >= 0 {
		name = path.Clean(name)
		if name[0] == '/' {
			name = name[1:]
		}
	}
	return name
}
