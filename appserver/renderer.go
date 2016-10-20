package main

import (
	"github.com/bughou-go/xiaomei/appserver/funcs"
	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/slice"
	"github.com/bughou-go/xm"
	"path"
	"strings"
)

func getRenderer() *xm.Renderer {
	var cache = config.Data.Env != `dev`
	var renderer_funcs = map[string]interface{}{
		`mt`:           funcs.AddModificationTimeFunc(cache),
		`html_safe`:    funcs.HtmlSafe,
		`dict`:         funcs.MakeDict,
		`keys`:         funcs.MapKeys,
		`values`:       funcs.MapValues,
		`IF`:           funcs.IF,
		`field`:        funcs.StructOrMapField,
		`union`:        slice.Union,
		`keys_union`:   funcs.MapKeysUnion,
		`thousand_sep`: utils.ThousandSep,
		`contains`:     strings.Contains,
	}
	return xm.NewRenderer(path.Join(config.Root, `views`), `layout/default`, cache, renderer_funcs)
}
