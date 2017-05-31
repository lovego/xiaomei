package renderer

import (
	// "fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"
	"strings"
)

type Renderer struct {
	Root      string
	Layout    string
	Templates map[string]*template.Template

	Funcs template.FuncMap
}

func New(root, layout string, cache bool, funcs template.FuncMap) *Renderer {
	r := &Renderer{
		Root:   path.Clean(root),
		Layout: path.Clean(layout),
		Funcs:  funcs,
	}
	if cache {
		r.Templates = make(map[string]*template.Template)
	}
	return r
}

func (r *Renderer) Render(wr io.Writer, name string, data interface{}, option O) {
	option = option.Process(r)
	tmpl := r.getTarget(name, option)
	var err error
	if option.HasLayout() {
		err = tmpl.ExecuteTemplate(wr, option.Layout, option.LayoutData(data))
	} else {
		err = tmpl.Execute(wr, data)
	}
	if err != nil {
		panic(err)
	}
}

func (r *Renderer) getTarget(name string, option O) *template.Template {
	name = cleanName(name)
	key := name
	if option.HasLayout() {
		key += `@` + option.Layout
	}
	tmpl := r.Templates[key]
	if tmpl == nil {
		if option.HasLayout() {
			tmpl = r.getTemplateWithLayout(name, option)
			if r.Templates != nil {
				r.Templates[key] = tmpl
			}
		} else {
			tmpl = r.getTemplate(name, name, option.Funcs)
		}
	}
	return tmpl
}
func (r *Renderer) getTemplateWithLayout(name string, option O) *template.Template {
	layoutTmpl, err := r.getTemplate(option.Layout, option.Layout, option.Funcs).Clone()
	if err != nil {
		panic(err)
	}
	for _, t := range r.getTemplate(name, ``, option.Funcs).Templates() {
		if _, err := layoutTmpl.AddParseTree(t.Name(), t.Tree); err != nil {
			panic(err)
		}
	}
	return layoutTmpl
}

// 针对同一个模板使用不同的funcs，因为缓存的原因，会得到非预期的结果。
func (r *Renderer) getTemplate(name, tmplName string, funcs template.FuncMap) *template.Template {
	tmpl := r.Templates[name]
	if tmpl == nil {
		tmpl = template.New(tmplName)
		if r.Funcs != nil {
			tmpl.Funcs(r.Funcs)
		}
		if funcs != nil {
			tmpl.Funcs(funcs)
		}
		parseTemplate(tmpl, name, r.Root, r.Root, make(map[string]bool))
		if r.Templates != nil {
			r.Templates[name] = tmpl
		}
	}
	return tmpl
}

// 关联模板中不能含有同名模板
func parseTemplate(tmpl *template.Template, name, base, root string, loaded map[string]bool) {
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
	template.Must(tmpl.Parse(string(content)))

	var new []string
	for _, t := range tmpl.Templates() { // 所有已关联的模板，包括自己
		nam := t.Name()
		if strings.IndexByte(nam, '/') >= 0 && !loaded[nam] {
			loaded[nam] = true
			if nam != name {
				new = append(new, nam)
			}
		}
	}

	base = path.Dir(p)
	for _, nam := range new {
		parseTemplate(tmpl.New(nam), nam, base, root, loaded)
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
