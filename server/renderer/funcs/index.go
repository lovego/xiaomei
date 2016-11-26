package funcs

import (
	"strings"

	"github.com/bughou-go/xiaomei/config"
	"github.com/bughou-go/xiaomei/utils"
	"github.com/bughou-go/xiaomei/utils/slice"
)

func Map() map[string]interface{} {
	return map[string]interface{}{
		`asset`:        AssetFunc(config.Data.Env == `dev`),
		`html_safe`:    HtmlSafe,
		`dict`:         MakeDict,
		`keys`:         MapKeys,
		`values`:       MapValues,
		`IF`:           IF,
		`field`:        StructOrMapField,
		`union`:        slice.Union,
		`keys_union`:   MapKeysUnion,
		`thousand_sep`: utils.ThousandSep,
		`contains`:     strings.Contains,
	}
}
