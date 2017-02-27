package funcs

import (
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils/number"
	"github.com/bughou-go/xiaomei/utils/slice"
)

func Map() map[string]interface{} {
	return map[string]interface{}{
		`asset`:        AssetFunc(config.Env() == `dev`),
		`html_safe`:    HtmlSafe,
		`dict`:         MakeDict,
		`keys`:         MapKeys,
		`values`:       MapValues,
		`IF`:           IF,
		`field`:        StructOrMapField,
		`union`:        slice.Union,
		`keys_union`:   MapKeysUnion,
		`thousand_sep`: number.ThousandSep,
		`contains`:     strings.Contains,
	}
}
